package main

import (
	"flag"
	"log"
	"time"

	"github.com/ujjwal-shekhar/mapreduce/services/common/user_code/template"
	"github.com/ujjwal-shekhar/mapreduce/services/server/handler/utils"
)

func main() {
	TYPE := flag.String("TYPE", "mapper", "Type of server to start (mapper = true, reducer = false)")
	PORT := flag.String("PORT", "5001", "Port to start the server on")
	NUM_REDUCERS := flag.Int("NUM_REDUCERS", 1, "Number of reducers to connect to")
	
	flag.Parse()
	log.Printf("Starting %s server on port %s", *TYPE, *PORT)
	time.Sleep(5 * time.Second)

	// Create a new worker server, and connect to all reducers
	workerServer := NewWorker(*TYPE == "mapper", "wordcount", *NUM_REDUCERS)
	if workerServer == nil {
		log.Fatalf("Failed to create worker server")
	}
	
	// Start the working process, collect the results, send to reducers
	outChannel := make(chan common.KV)
	go utils.StartWorking(workerServer, outChannel)
	if *TYPE == "mapper" {
		go utils.SendToReducers(workerServer, outChannel, PORT)
	} else {
		go utils.OutputFinalResults(outChannel, PORT)
	}

	StartWorker(*TYPE == "mapper", string(*PORT), "wordcount", workerServer)

	// log.Printf("Finished writing results to file %s", fileName)
}
