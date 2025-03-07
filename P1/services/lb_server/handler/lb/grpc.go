package handler

import (
	"context"
	"log"
	"sync"
	"time"
	"errors"

	"go.etcd.io/etcd/client/v3"

	pb "github.com/ujjwal-shekhar/load_balancer/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/load_balancer/services/lb_server/handler/utils"
	"github.com/ujjwal-shekhar/load_balancer/services/common/utils/constants"
)

// LoadBalancer struct
type LoadBalancer struct {
	pb.UnimplementedLoadBalancerServer
	etcdClient *clientv3.Client

	mu      sync.Mutex
	servers map[string]*pb.ServerInfo
	policy  utils.ServerSelectionPolicy
}

// NewLoadBalancer initializes the Load Balancer
func NewLoadBalancer(policyName string) *LoadBalancer {
	policyFunc := utils.GetPolicy(policyName)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.ETCD_ENDPOINT},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}

	lb := &LoadBalancer{
		servers:    make(map[string]*pb.ServerInfo),
		policy:     policyFunc,
		etcdClient: cli,
	}

	// Start watching etcd lease events
	go lb.watchEtcdLeases()

	return lb
}

// Watch etcd for lease events (renewals + expirations)
func (lb *LoadBalancer) watchEtcdLeases() {
	watchChan := lb.etcdClient.Watch(context.Background(), constants.ETCD_SERVERS_PREFIX, clientv3.WithPrefix())

	for watchResp := range watchChan {
		for _, ev := range watchResp.Events {
			serverAddr := string(ev.Kv.Key[len(constants.ETCD_SERVERS_PREFIX):])

			if ev.Type == clientv3.EventTypeDelete {
				// Lease expired -> Remove from servers
				lb.mu.Lock()
				delete(lb.servers, serverAddr)
				lb.mu.Unlock()
				log.Printf("Server %s removed (lease expired)", serverAddr)
			} else if ev.Type == clientv3.EventTypePut {
				log.Printf("Server %s lease renewed", serverAddr)
				lb.mu.Lock()
				lb.servers[serverAddr] = &pb.ServerInfo{Address: serverAddr}
				lb.mu.Unlock()

				// print all servers
				lb.mu.Lock()
				log.Println("Current servers:", len(lb.servers))
				for addr, server := range lb.servers {
					log.Printf("Server: %s, CPU Load: %f, Task Load: %d", addr, server.CpuLoad, server.TaskLoad)
				}
				lb.mu.Unlock()
			}

			// print event type
			log.Printf("Event Type: %v", ev.Type)
		}
	}
}

func (lb* LoadBalancer) ProcessServerHeartbeat(ctx context.Context, req *pb.ServerInfo) (*pb.ServerReply, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	serverAddr := req.Address
	lb.servers[serverAddr] = req

	log.Printf("Received heartbeat from server %s", serverAddr)
	return &pb.ServerReply{
		Success: true,
	}, nil
}

func (lb *LoadBalancer) ProcessClientRequest(ctx context.Context, req *pb.ClientRequest) (*pb.ServerInfo, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	server := lb.policy(lb.servers)
	if server == nil {
		return nil, errors.New("no servers available")
	}

	log.Printf("Selected server: %s", server.Address)
	return server, nil
}