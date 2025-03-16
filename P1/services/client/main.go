package main

import (
	"log"
	"os"
	"strconv"

	"github.com/ujjwal-shekhar/load_balancer/services/client/handler/client"
)

func main() {
	// Read the command line argument for task load amount
	taskLoad, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to parse task load: %v", err)
	}

	cli := handler.NewClient()
	cli.Load = int32(taskLoad)
	cli.Run()
}