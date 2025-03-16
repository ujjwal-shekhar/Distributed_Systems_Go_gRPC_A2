package main

import (
	"log"
	"net"
	// "fmt"
	// "os"

	"google.golang.org/grpc"

	pb "github.com/ujjwal-shekhar/mapreduce/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/mapreduce/services/common/user_code"
	"github.com/ujjwal-shekhar/mapreduce/services/server/handler/server"
)

func StartWorker(isMapper bool, port string, taskType string, workerServer *handler.Worker) {
	lis, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Worker gRPC server is running on port %v", lis.Addr().String())

	grpcServer := grpc.NewServer()
	
	log.Printf("Worker server created")
	pb.RegisterFileTransferServer(grpcServer, workerServer)

	log.Printf("Worker server registered")

	log.Printf("Worker gRPC server is running on port %s", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	log.Printf("Worker gRPC server is serving on port %s", lis.Addr().String())
}

func NewWorker(isMapper bool, taskType string, numReducers int) *handler.Worker {
	taskDetails := usercode.GetTaskDetails(taskType)
	if taskDetails == nil {
		log.Fatalf("Task details not found for task: %s", taskType)
	}
	log.Printf("Task details found for task: %s | %v", taskType, taskDetails)

	if isMapper {
		log.Printf("Creating new mapper")
		return handler.NewMapper(taskDetails, taskType, numReducers)
	} else {
		return handler.NewReducer(taskDetails, taskType, numReducers)
	}
}