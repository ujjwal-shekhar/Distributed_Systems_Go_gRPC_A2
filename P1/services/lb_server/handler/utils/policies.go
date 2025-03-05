package utils

import (
	"sync"
	"time"

	"github.com/ujjwal-shekhar/load_balancer/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/load_balancer/services/common/utils/constants"
)

var (
	roundRobinCounter 	string
 	roundRobinMutex  	sync.Mutex
)

// PickFirst always selects the first available server.
func PickFirst(servers map[string]*ServerMetadata) *load_balancer.ServerInfo {
	for _, server := range servers {
		// We check if the server is available which reduces to 
		// checking if the last received heartbeat is within fail threshold
		if time.Since(server.LastUpdated) < constants.SERVER_TIMEOUT {
			return server.Info
		}
	}
	return nil
}

// LeastLoaded selects the server with the lowest CPU load.
func LeastLoaded(servers map[string]*ServerMetadata) *load_balancer.ServerInfo {
	var bestServer *load_balancer.ServerInfo
	for _, server := range servers {
		if bestServer == nil {
			bestServer = server.Info
		} else if server.Info.CpuLoad < bestServer.CpuLoad && time.Since(server.LastUpdated) < constants.SERVER_TIMEOUT {
			bestServer = server.Info
		}
	}
	return bestServer
}

// RoundRobin selects servers in a cyclic manner.
func RoundRobin(servers map[string]*ServerMetadata) *load_balancer.ServerInfo {
	roundRobinMutex.Lock()
	defer roundRobinMutex.Unlock()

	if len(servers) == 0 {
		return nil
	}

	// Convert map to a slice of available servers
	activeServers := make([]*ServerMetadata, 0, len(servers))
	for _, server := range servers {
		if time.Since(server.LastUpdated) <= constants.SERVER_TIMEOUT {
			activeServers = append(activeServers, server)
		}
	}

	if len(activeServers) == 0 {
		return nil // No active servers
	}

	// Find the index of the last selected server
	candidate := -1
	for i, server := range activeServers {
		if server.Info.Address == roundRobinCounter {
			candidate = i
			break
		}
	}

	// If last selected server is not found, start from the first
	if candidate == -1 {
		candidate = 0
	} else {
		candidate = (candidate + 1) % len(activeServers) // Move to next server
	}

	// Ensure we don’t loop indefinitely
	startIndex := candidate
	for {
		server := activeServers[candidate]
		if time.Since(server.LastUpdated) <= constants.SERVER_TIMEOUT {
			roundRobinCounter = server.Info.Address
			return server.Info
		}

		// Move to next candidate
		candidate = (candidate + 1) % len(activeServers)

		// If we've wrapped around, return nil
		if candidate == startIndex {
			return nil
		}
	}
}