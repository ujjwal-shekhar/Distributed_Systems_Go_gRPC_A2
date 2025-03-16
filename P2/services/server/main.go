package main

import (
	"flag"
	"log"
	"strconv"
	// "time"
)

func main() {
	TYPE := flag.String("TYPE", "mapper", "Type of server to start (mapper = true, reducer = false)")
	PORT := flag.String("PORT", "5001", "Port to start the server on")
	NUM_REDUCERS := flag.Int("NUM_REDUCERS", 1, "Number of reducers to connect to")
	TASK := flag.String("TASK", "wordcount", "Task ID")
	flag.Parse()
	
	log.Printf("Starting %s server on port %s", *TYPE, *PORT)

	// Create a new worker server
	workerServer := NewWorker(*TYPE == "mapper", *TASK, *NUM_REDUCERS)
	if workerServer == nil {
		log.Fatalf("Failed to create worker server")
	}
	val, _ := strconv.Atoi(*PORT)
	workerServer.PortNumber = val
	
	log.Printf("Worker server numreduce : %d", workerServer.NumReducers)

	StartWorker(*TYPE == "mapper", string(*PORT), *TASK, workerServer)
}
