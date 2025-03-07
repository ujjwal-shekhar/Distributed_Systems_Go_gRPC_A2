package utils

import (
	"time"

	pb "github.com/ujjwal-shekhar/load_balancer/services/common/genproto/comms"
)


type ServerMetadata struct {
	Info        	*pb.ServerInfo
	LastUpdated 	time.Time
}

type ServerSelectionPolicy func(map[string]*pb.ServerInfo) *pb.ServerInfo
