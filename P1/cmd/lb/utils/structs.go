package utils

import (
	"time"
	"sync"

	pb "p1/protofiles/service"
	"go.etcd.io/etcd/client/v3"
)

const serverTTLTicks = 10
const serverTTLTicksDuration = time.Second * serverTTLTicks

type ServerInfo struct {
	Address   string
	CPULoad   float32
	TotalLoad int32
	LastPing  time.Time
}

type LoadBalancer struct {
	pb.UnimplementedLoadBalancerServer
	servers 	map[string]*ServerInfo
	mu      	sync.Mutex
	etcd    	*clientv3.Client
}