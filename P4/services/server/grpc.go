package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/ujjwal-shekhar/bft/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/bft/services/server/handler/server"
)

func StartServer(PORT int, server *server.Server) {
	// Start listening to the port
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", PORT))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening on port: %d", PORT)

	// Create a new server
	grpcServer := grpc.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	pb.RegisterOMServer(grpcServer, server)
	log.Printf("Server %d of type %s created", server.ID, server.TYPE)

	// Serve, NOW!
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	log.Printf("Serving on port: %d", PORT)

	select {}
}