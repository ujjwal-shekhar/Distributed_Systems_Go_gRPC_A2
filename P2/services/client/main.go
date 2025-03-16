package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/ujjwal-shekhar/mapreduce/services/client/handler/client"
	"github.com/ujjwal-shekhar/mapreduce/services/client/handler/utils"
	"github.com/ujjwal-shekhar/mapreduce/services/common/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	// We parse CLI args
	M := flag.Int("M", 5, "M : Number of mappers")
	R := flag.Int("R", 5, "R : Number of reducers")
	T := flag.String("T", "wordcount", "T : Task ID")
	flag.Parse()

	log.Printf("Mappers: %d, Reducers: %d, Task ID: %s", *M, *R, *T)

	// Lets create a master
	master := client.NewMaster(*M, *R, *T)
	log.Println("Master created")
	time.Sleep(5 * time.Second)


	// Begin the pipeline : master -> mappers
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := utils.FileToMappersPipeline(helpers.DATASET_PATH, master, &wg); err != nil {
			log.Fatalf("Error in FileToMappersPipeline: %v", err)
		}
	}()
	wg.Wait()
	time.Sleep(5 * time.Second)

	// Send all reducers the signal to start 
	// vomiting the outputs to secondary storage
	for _, reducer := range master.ReducerServers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := reducer.Vomit(context.Background(), &emptypb.Empty{})
			if err != nil {
				log.Fatalf("Error in Vomit: %v", err)
			}
		}()
	}
	wg.Wait()

	// Keep listening to reducers that inform about completion
	// then collect them here in a single output file
	log.Println("Client done")

	// // Close all connections
	// for _, conn := range master.MapperServers {
	// 	conn.Close()
	// }
	// for _, conn := range master.ReducerServers {
	// 	conn.Close()
	// }
	// select {}
}
