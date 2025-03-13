package gateway

import (
	"context"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	
	crm.mu.Lock()
	defer crm.mu.Unlock()
	if cachedResp, exists := crm.m[idempotencyKey[0]]; exists {
		log.Printf("IdempotencyKeyUnaryInterceptor: Found cached response for %s", idempotencyKey)
		return cachedResp.Resp, cachedResp.Err
	}

	// Process the request
	resp, err := handler(ctx, req)
	crm.m[idempotencyKey[0]] = CachedResponse{Resp: resp, Err: err}
	log.Printf("IdempotencyKeyUnaryInterceptor: Caching response for %s", idempotencyKey)
	return resp, err
}