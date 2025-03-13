package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/ujjwal-shekhar/stripe-clone/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/stripe-clone/services/gateway/handler/auth"
	"github.com/ujjwal-shekhar/stripe-clone/services/gateway/handler/gateway"
	"github.com/ujjwal-shekhar/stripe-clone/services/common/utils"
)

func StartGatewayServer(reqPath string, keyPath string) {
	tlsConfig, err := auth.LoadTLSCredentials(reqPath, keyPath)
	if err != nil {
		log.Fatalf("failed to load TLS credentials: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsConfig), 
		grpc.UnaryInterceptor(
			auth.RBACUnaryInterceptor,
		),
		// grpc.WithChainUnaryInterceptor(
		// 	auth.RBACUnaryInterceptor,
		// 	auth.LoggerUnaryInterceptor,
		// ),
	)

	gatewayServer := gateway.NewGatewayServerTLS(tlsConfig)
	pb.RegisterStripeServiceServer(grpcServer, gatewayServer)

	lis, err := net.Listen("tcp", utils.PAYMENT_GATEWAY_URL)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Payment gateway gRPC server is running on port %s", utils.PAYMENT_GATEWAY_URL)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}