package auth

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryInterceptor to inject JWT token into metadata
func TokenUnaryInterceptor(token string, ) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Attach token to metadata
		md := metadata.Pairs("authorization", "Bearer "+token)
		ctx = metadata.NewOutgoingContext(ctx, md)

		log.Printf("TokenUnaryInterceptor: token=%v md=%v", token, md)

		// Invoke the RPC with the updated context
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
