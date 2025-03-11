package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/ujjwal-shekhar/mapreduce/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/mapreduce/services/server/handler/server"
)

func StartWorker(isMapper bool, port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	workerServer := NewWorker(isMapper)

	pb.RegisterFileTransferServer(grpcServer, workerServer)
	reflection.Register(grpcServer)

	log.Printf("Worker gRPC server is running on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func NewWorker(isMapper bool) *handler.Worker {
	if isMapper {
		return handler.NewMapper()
	} else {
		return handler.NewReducer()
	}
}