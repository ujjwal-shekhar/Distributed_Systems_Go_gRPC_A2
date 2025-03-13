package auth

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

// LoggerUnaryInterceptor logs details of all incoming unary gRPC requests.
func LoggerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	log.Printf("[gRPC] Method: %s - Started", info.FullMethod)

	// Process request
	resp, err := handler(ctx, req)

	// Log completion with execution time and error (if any)
	duration := time.Since(start)
	if err != nil {
		log.Printf("[gRPC] Method: %s - Failed (%s) - Error: %v", info.FullMethod, duration, err)
	} else {
		log.Printf("[gRPC] Method: %s - Completed (%s)", info.FullMethod, duration)
	}

	return resp, err
}
