package utils

import (
	"sync"

	pb "github.com/ujjwal-shekhar/load_balancer/services/common/genproto/comms"
)

var (
	lastServerAddr string
	lastServerMu   sync.Mutex
)

// Meh just put the selector here to reduce clutter from the gRPC logic in handler
func GetPolicy(name string) ServerSelectionPolicy {
	switch name {
	case "pick_first":
		return PickFirst
	case "round_robin":
		return RoundRobin
	default:
		return LeastLoaded
	}
}

// PickFirst always selects the first available server.
func PickFirst(servers map[string]*pb.ServerInfo) *pb.ServerInfo {
	for _, server := range servers {
		return server
	}

	return nil
}

// LeastLoaded selects the server with the lowest CPU load.
func LeastLoaded(servers map[string]*pb.ServerInfo) *pb.ServerInfo {
	bestServer := &pb.ServerInfo{CpuLoad: 100, TaskLoad: 1000000}
	for _, server := range servers {
		if server.CpuLoad < bestServer.CpuLoad {
			bestServer = server
		} else if server.CpuLoad == bestServer.CpuLoad && server.TaskLoad < bestServer.TaskLoad {
			bestServer = server
		}
	}

	return bestServer
}

// RoundRobin selects servers in a cyclic manner.
func RoundRobin(servers map[string]*pb.ServerInfo) *pb.ServerInfo {
	lastServerMu.Lock()
	defer lastServerMu.Unlock()

	if len(servers) == 0 {
		return nil
	}

	found_last := len(servers) == 1
	for addr, server := range servers {
		if found_last {
			lastServerAddr = addr
			return server
		}

		found_last = addr == lastServerAddr
	}

	for addr, server := range servers {
		lastServerAddr = addr
		return server
	}

	return nil
}