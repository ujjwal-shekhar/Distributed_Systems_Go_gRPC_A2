package main

import (
	"github.com/ujjwal-shekhar/bft/services/client/handler/client"
)

func StartSimulation(N int, T int) {
	client.ForkProcs(N, T)
}
