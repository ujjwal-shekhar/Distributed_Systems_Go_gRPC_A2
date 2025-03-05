package utils

import (
	"time"

	"github.com/ujjwal-shekhar/load_balancer/services/common/genproto/comms"
)


type ServerMetadata struct {
	Info        	*load_balancer.ServerInfo
	LastUpdated 	time.Time
}

type ServerSelectionPolicy func(map[string]*ServerMetadata) *load_balancer.ServerInfo
