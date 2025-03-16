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
	var errChan = make(chan error, numMappers+numReducers) // Channel to collect errors

	// Spawn and connect to mappers
	for i := 0; i < numMappers; i++ {
		port := 5000 + i
		wg.Add(1)
		go func(ii int, port int) {
			defer wg.Done()

			// Spawn process
			cmd := exec.Command("make", "run-server", "TYPE=mapper", fmt.Sprintf("PORT=%d", port))
			if err := cmd.Start(); err != nil {
				errChan <- fmt.Errorf("failed to start mapper process on port %d: %v", port, err)
				return
			}

			// Wait for the server to start (retry mechanism)
			addr := fmt.Sprintf("localhost:%d", port)
			var conn *grpc.ClientConn
			var err error
			for retry := 0; retry < 5; retry++ {
				conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err == nil {
					break
				}
				time.Sleep(500 * time.Millisecond) // Wait before retrying
			}
			if err != nil {
				errChan <- fmt.Errorf("failed to connect to mapper at %s: %v", addr, err)
				return
			}

			// Store connection
			master.MapperServers[ii] = pb.NewFileTransferClient(conn)
			log.Printf("Connected to mapper at %s", addr)
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
			if err := cmd.Start(); err != nil {
				errChan <- fmt.Errorf("failed to start reducer process on port %d: %v", port, err)
				return
			}

			// Wait for the server to start (retry mechanism)
			addr := fmt.Sprintf("localhost:%d", port)
			var conn *grpc.ClientConn
			var err error
			for retry := 0; retry < 5; retry++ {
				conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err == nil {
					break
				}
				time.Sleep(500 * time.Millisecond) // Wait before retrying
			}
			if err != nil {
				errChan <- fmt.Errorf("failed to connect to reducer at %s: %v", addr, err)
				return
			}

			// Store connection
			master.ReducerServers[i] = pb.NewFileTransferClient(conn)
			log.Printf("Connected to reducer at %s", addr)
		}(i, port)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Check for errors
	select {
	case err := <-errChan:
		log.Fatalf("Error during process spawning or connection: %v", err)
	default:
		// No errors, continue
	}

	return master
}