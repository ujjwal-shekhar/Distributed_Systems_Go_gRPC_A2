package main

import (
	"context"
	"log"
	"time"

	pb "github.com/ujjwal-shekhar/load_balancer/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/load_balancer/services/common/utils/constants"
	"google.golang.org/grpc"
)

// Server represents a backend server that registers with the load balancer.
type Server struct {
	pb.LoadBalancerClient

	address   	string
	cpuLoad   	float32
	taskLoad 	int32
}

// NewServer initializes a new backend server.
func NewServer(lbAddr string, serverAddr string) *Server {
	// Dial the load balancer
	conn, err := grpc.NewClient(lbAddr)
	if err != nil {
		log.Fatalf("Failed to connect to load balancer: %v", err)
	}
	lbClient := pb.NewLoadBalancerClient(conn)

	// Return the server instance
	return &Server{
		LoadBalancerClient: lbClient,
		address:   serverAddr,
		cpuLoad:   0.0, // TODO : Use gopsutil here
		taskLoad:  0,
	}
}

// SendHeartbeat sends periodic heartbeats to the load balancer.
func (s *Server) SendHeartbeat() {
	ticker := time.NewTicker(constants.HEARTBEAT_INTERVAL * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		_, err := s.ProcessServerHeartbeat(context.Background(), &pb.ServerInfo{
			Address:  s.address,
			CpuLoad:  s.cpuLoad,
			TaskLoad: s.taskLoad,
		})

		if err != nil {
			log.Printf("Failed to send heartbeat: %v", err)
		}
	}
}
