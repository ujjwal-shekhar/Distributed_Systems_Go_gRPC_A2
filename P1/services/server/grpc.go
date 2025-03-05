package main

import (
	"log"
	"net"

	pb "github.com/ujjwal-shekhar/load_balancer/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/load_balancer/services/common/utils/constants"
)

// Have a go routine to send heartbeats to the lb


func main() {
	// Create a new server instance
	server := NewServer(constants.LB_ADDRESS, constants.SERVER_ADDRESS)
