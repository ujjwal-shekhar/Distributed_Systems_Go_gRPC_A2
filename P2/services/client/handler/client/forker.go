package client

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/ujjwal-shekhar/mapreduce/services/common/genproto/comms"
)

func ForkProcs(numMappers int, numReducers int) *Master {
	master := &Master{
		MapperServers:  make([]pb.FileTransferClient, numMappers),
		ReducerServers: make([]pb.FileTransferClient, numReducers),
	}

	var wg sync.WaitGroup

	// Spawn and connect to mappers
	for i := 0; i < numMappers; i++ {
		port := 5000 + i
		wg.Add(1)
		go func(i, port int) {
			defer wg.Done()

			// Spawn process
			cmd := exec.Command("make", "run-server", "TYPE=mapper", fmt.Sprintf("PORT=%d", port))
			cmd.Start()
			time.Sleep(500 * time.Millisecond)

			// Dial gRPC connection
			addr := fmt.Sprintf("localhost:%d", port)
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Failed to connect to mapper at %s: %v", addr, err)
			}

			// Store connection
			master.MapperServers[i] = pb.NewFileTransferClient(conn)
		}(i, port)
	}

	// Spawn and connect to reducers
	for i := 0; i < numReducers; i++ {
		port := 6000 + i
		wg.Add(1)
		go func(i, port int) {
			defer wg.Done()

			// Spawn process
			cmd := exec.Command("make", "run-server", "TYPE=reducer", fmt.Sprintf("PORT=%d", port))
			cmd.Start()
			time.Sleep(500 * time.Millisecond)

			addr := fmt.Sprintf("localhost:%d", port)
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Failed to connect to reducer at %s: %v", addr, err)
			}

			// Store connection
			master.ReducerServers[i] = pb.NewFileTransferClient(conn)
		}(i, port)
	}

	// Wait for all connections to complete
	wg.Wait()

	return master
}
