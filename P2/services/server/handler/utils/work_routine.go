package utils

import (
	// "fmt"
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"strconv"

	// "os"

	pb "github.com/ujjwal-shekhar/mapreduce/services/common/genproto/comms"
	// task1 "github.com/ujjwal-shekhar/mapreduce/services/common/user_code/Task1"
	"github.com/ujjwal-shekhar/mapreduce/services/common/user_code/template"
	handler "github.com/ujjwal-shekhar/mapreduce/services/server/handler/server"
)

func hashKey(key string) int {
	// Hash the key
	hasher := md5.New()
	hasher.Write([]byte(key))
	hash := hasher.Sum(nil)
	hasher.Reset()

	// Find the reducer to send to
	return int(hash[0])
}

func StartWorking(worker *handler.Worker, outChannel chan common.KV) {
	log.Printf("Starting worker routine!!!!!!!!!!")
	for kv := range worker.TaskChannel {
		log.Printf("Received key-value pair")
		if worker.TaskDesc == "mapper" {
			worker.WorkFunc([]common.KV{kv}, outChannel)
		} else {

		}
	}
}

func SendToReducers(worker *handler.Worker, outChannel chan common.KV, PORT *string) {
	var Buffer []any

		// We will push things on to the buffer until
		// it reaches a size of 10000 key-value pairs
		// Then we will send it 
	for kv := range outChannel {
		// Push the key-value pair to the buffer
		Buffer = append(Buffer, kv)
		if len(Buffer) == 10000 {
			// Send the buffer to the reducer
			MakeNewFile(Buffer, PORT)
			Buffer = nil
		}
	}
}

func MakeNewFile(Buffer []any, PORT *string) {
	

func OutputFinalResults(outChannel chan common.KV, PORT *string) {
	// Open a file with name <port>.out
		// // Open a file with name <port>.out
	fileName := fmt.Sprintf("tmp/%s.out", *PORT)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", fileName, err)
	}
	defer file.Close()

	// Write everything from outChannel to the file
	for kv := range outChannel {
		// Convert the key-value pair to a string
		line := fmt.Sprintf("Key: %v, Value: %v\n", kv.GetKey(), kv.GetValue())

		// Write the line to the file
		_, err := file.WriteString(line)
		if err != nil {
			log.Printf("Failed to write to file %s: %v", fileName, err)
			continue
		}

		// Flush the file to ensure the data is written
		err = file.Sync()
		if err != nil {
			log.Printf("Failed to flush file %s: %v", fileName, err)
		}
	}
}