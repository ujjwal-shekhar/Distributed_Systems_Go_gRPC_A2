package client

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"sync"

	"github.com/ujjwal-shekhar/bft/services/common/utils"
)

func ForkProcs(N int, T int) {
	traitors := getTraitors(N, T)

	var wg sync.WaitGroup
	var errChan = make(chan error, N) // Channel to collect errors

	// Spawn and connect to others
	for id, general_type := range traitors {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()

			// Spawn process
			cmd := exec.Command(
				"make", "run-server",
				fmt.Sprintf("ID=%d", ii),
				fmt.Sprintf("PORT=%d", utils.PORT_BASE+ii),
				fmt.Sprintf("TYPE=%s", general_type),
				fmt.Sprintf("N=%d", N),
				fmt.Sprintf("T=%d", T),
			)

			if err := cmd.Start(); err != nil {
				errChan <- fmt.Errorf("failed to start process on port %d: %v", 5000+ii, err)
				return
			}
			
			log.Printf("Started process on port : %d", utils.PORT_BASE+ii)
		}(id)
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
}

func getTraitors(N int, T int) map[int]string {
	// Pick any T random integers without replacement from [1..N]
	// This will be the set of traitors
	traitors := make(map[int]string)
	for len(traitors) < T {
		traitors[rand.Intn(N)+1] = "traitor"
	}

	// Makee all others "honest"
	for i := range N {
		if _, ok := traitors[i]; !ok {
			traitors[i] = "honest"
		}
	}

	return traitors
}