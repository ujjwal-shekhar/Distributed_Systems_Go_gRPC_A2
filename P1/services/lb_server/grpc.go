package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/ujjwal-shekhar/load_balancer/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/load_balancer/services/common/utils/constants"
	"github.com/ujjwal-shekhar/load_balancer/services/lb_server/handler/lb"
)

func StartLBServer() {
	lis, err := net.Listen("tcp", constants.LB_PORT)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	loadBalancer := handler.NewLoadBalancer("least_loaded")
	pb.RegisterLoadBalancerServer(grpcServer, loadBalancer)

	reflection.Register(grpcServer)

	log.Printf("Load balancer gRPC server is running on port %s", constants.LB_PORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
