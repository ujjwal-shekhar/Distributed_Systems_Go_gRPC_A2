package types

import (
	"context"

	"github.com/ujjwal-shekhar/load_balancer/services/common/genproto/comms"
)

type LoadBalancerServer interface {
	ProcessClientRequest(ctx context.Context, req *load_balancer.ClientRequest) (*load_balancer.ServerInfo, error)
	ProcessServerHeartbeat(ctx context.Context, info *load_balancer.ServerInfo) (*load_balancer.ServerReply, error)
}
