package utils

import (
	// "sync"

	// pb "p1/protofiles/service"
)

type RoundRobinPolicy struct {
	last_taken int32
}

func (p *RoundRobinPolicy) SelectServer(lb *LoadBalancer) *ServerInfo {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.servers) == 0 {
		return nil
	}



	// keys := make([]string, 0, len(lb.servers))
	// for k := range lb.servers {
	// 	keys = append(keys, k)
	// }
	// server := lb.servers[keys[p.counter]]

	// p.counter++
	// p.counter %= int32(len(lb.servers))

	return server
}

type LeastLoadPolicy struct {}

func (p *LeastLoadPolicy) SelectServer(lb *LoadBalancer) *ServerInfo {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.servers) == 0 {
		return nil
	}

	var bestServer *ServerInfo
	for _, server := range lb.servers {
		if bestServer == nil || server.TotalLoad < bestServer.TotalLoad {
			bestServer = server
		}
	}

	return bestServer
}