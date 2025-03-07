package main

import (
	"log"

	pb "github.com/ujjwal-shekhar/load_balancer/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/load_balancer/services/server/handler/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Have a go routine to send heartbeats to the lb

func StartTaskServer() {
	grpcServer := grpc.NewServer()
	taskServer, lis := handler.NewServer()
	pb.RegisterTaskRunnerServer(grpcServer, taskServer)

	reflection.Register(grpcServer)

	log.Printf("Task server listening on %s", taskServer.Address)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	go taskServer.SendHeartbeats()
}