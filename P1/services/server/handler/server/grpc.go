package handler

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"go.etcd.io/etcd/client/v3"

	pb "github.com/ujjwal-shekhar/load_balancer/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/load_balancer/services/common/utils/constants"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	pb.UnimplementedTaskRunnerServer
	etcdClient *clientv3.Client
	lbClient   pb.LoadBalancerClient

	Address    string
	TaskLoad   int32
	leaseID    clientv3.LeaseID
}

// NewServer initializes a new backend server.
func NewServer() (*Server, net.Listener) {
	// Get a random port for the server
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Make this my address in info
	serverAddr := fmt.Sprintf("localhost:%d", lis.Addr().(*net.TCPAddr).Port)

	// Connect to etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.ETCD_ENDPOINT},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}

	// Register with etcd
	leaseResp, err := cli.Grant(context.Background(), constants.TTL)
	if err != nil {
		log.Fatalf("Failed to grant lease: %v", err)
	}

	_, err = cli.Put(context.Background(), constants.ETCD_SERVERS_PREFIX+serverAddr, serverAddr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		log.Fatalf("Failed to register with etcd: %v", err)
	}

	// Connect to the load balancer
	conn, err := grpc.NewClient(constants.LB_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to load balancer: %v", err)
	}

	lbClient := pb.NewLoadBalancerClient(conn)

	go func() {
		keepAliveChan, err := cli.KeepAlive(context.Background(), leaseResp.ID)
		if err != nil {
			log.Fatalf("Failed to keep lease alive: %v", err)
		}
	
		for ka := range keepAliveChan {
			if ka == nil {
				log.Println("Lease expired!")
				return
			}
			log.Printf("Lease renewed for server %s", serverAddr)
		}
	}()

	// Start the server
	return &Server{
		etcdClient: cli,
		lbClient:   lbClient,
		Address:    serverAddr,
		TaskLoad:   0,
		leaseID:    leaseResp.ID,
	}, lis
}

// SendHeartbeat sends periodic heartbeats to the load balancer.
func (s *Server) SendHeartbeats() {
	ticker := time.NewTicker(constants.HEARTBEAT_INTERVAL * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		_, err := s.lbClient.ProcessServerHeartbeat(context.Background(), &pb.ServerInfo{
			Address:  s.Address,
			CpuLoad:  s.CpuLoad,
			TaskLoad: s.TaskLoad,
		})

		if err != nil {
			log.Printf("Failed to send heartbeat: %v", err)
		}
	}
}

func (s *Server) RunTask(ctx context.Context, req *pb.ClientRequest) (*pb.ServerReply, error) {
	s.TaskLoad++
	s.
	return &pb.ServerReply{Address: s.Address}, nil
}