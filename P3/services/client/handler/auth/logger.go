package auth

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

// LoggerUnaryInterceptor logs details of all outgoing unary gRPC requests.
func LoggerUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		log.Printf("[gRPC] Method: %s - Started", method)

		// Invoke the RPC with the updated context
		err := invoker(ctx, method, req, reply, cc, opts...)

		// Log completion with execution time and error (if any)
		duration := time.Since(start)
		if err != nil {
			log.Printf("[gRPC] Method: %s - Failed in (%s) - Error: %v", method, duration, err)
		} else {
			log.Printf("[gRPC] Method: %s - Completed in (%s)", method, duration)
		}

		return err
	}
}