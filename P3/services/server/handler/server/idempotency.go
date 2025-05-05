package bank

import (
	"context"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type CachedResponse struct {
	Resp interface{}
	Err  error
}

type CachedResponseMap struct {
	mu sync.Mutex
	m  map[string]CachedResponse
}

func NewCachedResponseMap() *CachedResponseMap {
	return &CachedResponseMap{
		m: make(map[string]CachedResponse),
	}
}

func (crm *CachedResponseMap) IdempotencyKeyUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Check if the request is already cached
	md, ok := metadata.FromIncomingContext(ctx); if !ok {
		return handler(ctx, req)
	}
	idempotencyKey, exists := md["idempotency-key"]; if !exists {
		return handler(ctx, req)
	}
	log.Printf("IdempotencyKeyUnaryInterceptor: Idempotency key found: %s", idempotencyKey)
	
	crm.mu.Lock()////////
	if cachedResp, exists := crm.m[idempotencyKey[0]]; exists {
		log.Printf("IdempotencyKeyUnaryInterceptor: Found cached response for %s", idempotencyKey)
		crm.mu.Unlock()//////
		return cachedResp.Resp, cachedResp.Err
	}
	crm.mu.Unlock()//////
	
	// Process the request
	resp, err := handler(ctx, req)

	crm.mu.Lock()////////
	if (status.Code(err) != codes.Unavailable) {
		crm.m[idempotencyKey[0]] = CachedResponse{Resp: resp, Err: err}
	}
	crm.mu.Unlock()//////

	log.Printf("IdempotencyKeyUnaryInterceptor: Caching response for %s", idempotencyKey)
	return resp, err
}