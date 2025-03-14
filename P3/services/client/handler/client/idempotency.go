package client

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// PendingTransaction stores information for retrying requests.
// IMP: We do not want to redo IdempotencyKeyUnaryInterceptor
// so the entire state is preserved in this struct.
type PendingTransaction struct {
	invoker 	grpc.UnaryInvoker
	ctx     	context.Context
	method 	  	string
	req     	interface{}
	reply   	interface{}
	cc       	*grpc.ClientConn
	opts     	[]grpc.CallOption
	lastTime 	time.Time
	resultCh 	chan error
	numRetries 	int
	mu 		   	sync.Mutex
}

// TransactionManager manages retries for failed transactions.
type TransactionManager struct {
	rwLock       sync.Mutex 
	pendingQueue map[string]*PendingTransaction
	threshold    time.Duration
}

// NewTransactionManager creates a new instance and starts retry loop.
func NewTransactionManager() *TransactionManager {
	tm := &TransactionManager{
		pendingQueue: make(map[string]*PendingTransaction),
		threshold:    RETRY_THRESHOLD,
	}

	// Start the background retry routine
	go tm.retryPendingTransactions()
	return tm
}

// IdempotencyKeyUnaryInterceptor generates an idempotency key and tracks requests.
func (tm *TransactionManager) IdempotencyKeyUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Attach the key to metadata
		// Append, dont overwrite
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.Pairs("idempotency-key", uuid.New().String())
		} else {
			md.Append("idempotency-key", uuid.New().String())
		}
		ctx = metadata.NewOutgoingContext(ctx, md)
		log.Printf("Generated idempotency key: %s", md["idempotency-key"][0])

		// The wait chan
		resultCh := make(chan error, 1)

		// Store the request in pending queue
		tm.rwLock.Lock() // Exclusive write access
		tm.pendingQueue[md["idempotency-key"][0]] = &PendingTransaction{
			invoker:  invoker,
			ctx:      ctx,
			method:   method,
			req:      req,
			reply:    reply,
			cc:       cc,
			opts:     opts,
			lastTime: time.Now(),
			resultCh: resultCh,
			numRetries: MAX_RETRIES,
		}
		tm.rwLock.Unlock()

		log.Printf("IdempotencyKeyUnaryInterceptor: md=%s", md)
		// I am ALLOWING a race condition to happen 
		// between the normal invoker and the retry routine
		
		// Invoke the gRPC call in a goroutine
		go func() {
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				// If the initial call succeeds, send the result and clean up
				resultCh <- nil
				tm.rwLock.Lock()
				delete(tm.pendingQueue, md["idempotency-key"][0])
				tm.rwLock.Unlock()
			} else {
				// If the initial call fails, let the retry routine handle it
				log.Printf("Initial request failed for key: %s, error: %v. Waiting for retry...", md["idempotency-key"][0], err)
			}
		}()

		// Wait for either the initial call or the retry to complete
		select {
		case err := <-resultCh:
			// Initial call or retry succeeded/failed
			return err
		case <-ctx.Done():
			// Context expired before retry could complete
			log.Printf("Context expired for key: %s, error: %v", md["idempotency-key"][0], ctx.Err())
			return ctx.Err()
		}
	}
}

// retryPendingTransactions retries failed requests after the threshold is exceeded.
func (tm *TransactionManager) retryPendingTransactions() {
	for {
		time.Sleep(RETRY_FREQUENCY)
		now := time.Now()

		// Lock for reading the pendingQueue
		tm.rwLock.Lock()
		for key, txn := range tm.pendingQueue {
			if now.Sub(txn.lastTime) > tm.threshold {
				log.Printf("Retrying request for key: %s", key)

				// Launch a goroutine for non-blocking retry
				go func(key string, txn *PendingTransaction) {
					// Lock the transaction for thread-safe access
					txn.mu.Lock()
					txn.numRetries--
					txn.lastTime = now // Update last retry time
					// Decrement retry count
					if txn.numRetries == 0 {
						log.Printf("Max retries reached for key: %s", key)
						txn.resultCh <-errors.New("max retries reached")
					}
					txn.mu.Unlock()

					// Retry request without interceptor to avoid adding a new idempotency key
					err := txn.invoker(txn.ctx, txn.method, txn.req, txn.reply, txn.cc, txn.opts...)
					if err == nil {
						log.Printf("Retry succeeded for key: %s", key)
						txn.resultCh <- nil
					} else {
						log.Printf("Retry failed for key: %s, error: %v", key, err)
						txn.resultCh <- err
					}

					// Remove the transaction from the pendingQueue
					tm.rwLock.Lock()
					delete(tm.pendingQueue, key)
					tm.rwLock.Unlock()
				}(key, txn)
			}
		}
		tm.rwLock.Unlock()
	}
}