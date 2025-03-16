package main

import (
	"flag"
	"log"
	"time"

	"github.com/ujjwal-shekhar/mapreduce/services/client/handler/client"
	"github.com/ujjwal-shekhar/mapreduce/services/client/handler/utils"
	"github.com/ujjwal-shekhar/mapreduce/services/common/utils"
)

func main() {
	// We parse CLI args
	M := flag.Int("M", 5, "M : Number of mappers")
	R := flag.Int("R", 5, "R : Number of reducers")
	T := flag.String("T", "wordcount", "T : Task ID")
	flag.Parse()

	log.Printf("Mappers: %d, Reducers: %d, Task ID: %s", *M, *R, *T)

	// Lets create a master
	master := client.NewMaster(*M, *R)
	log.Println("Master created")
	time.Sleep(10 * time.Second)

	// Begin the pipeline : chunk -> mappers
	chunkChannel := make(chan utils.ChunkMetadata)
	go utils.FileToMappersPipeline(helpers.DATASET_PATH, chunkChannel) // Producer
	go utils.SendChunksToMappers(chunkChannel, master) // Consumer

	// Keep listening to reducers that inform about completion
	// then collect them here in a single output file
	log.Println("Client done")

	select {} // TODO: Add a graceful shutdown mechanism
}
