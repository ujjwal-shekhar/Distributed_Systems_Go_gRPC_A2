package client

import (
	pb "github.com/ujjwal-shekhar/mapreduce/services/common/genproto/comms"
)

type Master struct {
	MapperServers 	[]pb.FileTransferClient
	ReducerServers 	[]pb.FileTransferClient
}

func NewMaster(numMappers int, numReducers int) *Master {
	return ForkProcs(numMappers, numReducers)
}
