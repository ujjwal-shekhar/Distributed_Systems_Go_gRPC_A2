package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "p1/protofiles/service"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	"p1/cmd/lb/utils"
)

var policy = "round-robin"

// Register or update server heartbeat
func (lb *LoadBalancer) ProcessServerHeartbeat(ctx context.Context, info *pb.ServerInfo) (*emptypb.Empty, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.servers[info.Address] = &ServerInfo{
		Address:   info.Address,
		CPULoad:   info.CpuLoad,
		TotalLoad: info.TotalLoad,
		LastPing:  time.Now(),
	}

	// Store in etcd with a lease (expiry)
	leaseResp, err := lb.etcd.Grant(ctx, serverTTLTicks) // 10-sec lease
	if err != nil {
		return nil, err
	}
	_, err = lb.etcd.Put(ctx, "servers/"+info.Address, fmt.Sprintf("%f,%d", info.CpuLoad, info.TotalLoad), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return nil, err
	}

	log.Printf("Server updated: %s | CPU: %.2f | Load: %d", info.Address, info.CpuLoad, info.TotalLoad)
	return &emptypb.Empty{}, nil
}

func (lb *LoadBalancer) ProcessClientRequest(ctx context.Context, req *pb.ClientRequest) (*pb.ServerInfo, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.servers) == 0 {
		return nil, fmt.Errorf("no servers available")
	}

	// Find server with least total load
	var bestServer *ServerInfo
	for _, server := range lb.servers {
		if bestServer == nil || server.TotalLoad < bestServer.TotalLoad {
			bestServer = server
		}
	}

	if bestServer == nil {
		return nil, fmt.Errorf("no servers found")
	}

	log.Printf("Routing client to server: %s", bestServer.Address)
	return &pb.ServerInfo{
		Address:   bestServer.Address,
		CpuLoad:   bestServer.CPULoad,
		TotalLoad: bestServer.TotalLoad,
	}, nil
}

func main() {
	// Setup etcd connection
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer etcdClient.Close()

	lb := &LoadBalancer{
		servers: make(map[string]*ServerInfo),
		etcd:    etcdClient,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLoadBalancerServiceServer(grpcServer, lb)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Load balancer started on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
