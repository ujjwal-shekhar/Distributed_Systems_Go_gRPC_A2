package handler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ujjwal-shekhar/load_balancer/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/load_balancer/services/lb_server/handler/utils"
)

// LoadBalancer struct with a dynamic policy function
type LoadBalancer struct {
	load_balancer.UnimplementedLoadBalancerServer

	mu      			sync.Mutex
	servers 			map[string]*utils.ServerMetadata
	policy  			utils.ServerSelectionPolicy
}

// NewLoadBalancer creates a new LoadBalancer instance with a given policy.
func NewLoadBalancer(policyName string) *LoadBalancer {
	var policyFunc utils.ServerSelectionPolicy
	switch policyName {
	case "pick_first":
		policyFunc = utils.PickFirst
	case "least_loaded":
		policyFunc = utils.LeastLoaded
	case "round_robin":
		policyFunc = utils.RoundRobin
	default:
		policyFunc = utils.LeastLoaded
	}

	return &LoadBalancer{
		servers: make(map[string]*utils.ServerMetadata),
		policy:  policyFunc,
	}
}

// ProcessClientRequest selects the best server based on the configured policy.
func (lb *LoadBalancer) ProcessClientRequest(ctx context.Context, req *load_balancer.ClientRequest) (*load_balancer.ServerInfo, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	bestServer := lb.policy(lb.servers) // Use the selected policy function
	if bestServer == nil {
		return nil, fmt.Errorf("no available servers")
	}

	return bestServer, nil
}

// ProcessServerHeartbeat updates the server metadata.
func (lb* LoadBalancer) ProcessServerHeartbeat (ctx context.Context, req *load_balancer.ServerInfo) (*load_balancer.ServerReply, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.servers[req.Address] = &utils.ServerMetadata{
		Info:        req,
		LastUpdated: time.Now(),
	}

	rep := &load_balancer.ServerReply{
		Success: true,
	}

	log.Printf("Received heartbeat from server: %s", req.Address)

	return rep, nil
}
