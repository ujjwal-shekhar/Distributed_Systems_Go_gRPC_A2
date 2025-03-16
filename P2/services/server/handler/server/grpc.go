package handler

import (
	"context"
	"fmt"
	"log"
	"sync"

	pb "github.com/ujjwal-shekhar/mapreduce/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/mapreduce/services/common/user_code"
	"github.com/ujjwal-shekhar/mapreduce/services/common/user_code/template"
	"github.com/ujjwal-shekhar/mapreduce/services/server/handler/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Worker struct {
	pb.UnimplementedFileTransferServer
	
	TaskChannel 	chan common.KV
	WorkFunc 		func([]common.KV, chan common.KV)
	UserDetails     *usercode.UserTaskDetails
	
	mu 				sync.Mutex
	ReducerList		[]utils.KV
	IsMapper		bool
	TaskDesc		string
	NumReducers		int
	PortNumber		int
}

func NewMapper(taskDetails *usercode.UserTaskDetails, taskType string, numReducers int) *Worker {
	return &Worker {
		mu: sync.Mutex{},
		IsMapper: true,
		TaskDesc: taskType,
		TaskChannel: make(chan common.KV),
		WorkFunc: taskDetails.Mapper,
		UserDetails: taskDetails,
		NumReducers: numReducers,
	}
}

func NewReducer(taskDetails *usercode.UserTaskDetails, taskType string, numReducers int) *Worker {
	return &Worker{
		mu: sync.Mutex{},
		IsMapper: false,
		TaskDesc: taskType,
		TaskChannel: make(chan common.KV),
		WorkFunc: taskDetails.Reducer,
		UserDetails: taskDetails,
		NumReducers: numReducers,
	}
}

func (w *Worker) SendToMapper(stream pb.FileTransfer_SendToMapperServer) error {
	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("Finished receiving chunks.")
				break
			}
			log.Printf("Error receiving chunk: %v", err)
			return err
		}

		// Apply the map function according to the task
		switch w.TaskDesc {
		case "wordcount":
			// Split the chunk into words
			words := utils.SplitChunksToWords(chunk.Chunk)
			for _, word := range words {
				// Emit the key-value pair (word, "1") into the intermediate file
				utils.EmitKeyValue(word, "1", w.PortNumber, w.NumReducers)
			}

		case "invertedindex":
			// Split the chunk into words
			words := utils.SplitChunksToWords(chunk.Chunk)
			for _, word := range words {
				// Emit the key-value pair (word, chunk.FileName) into the intermediate file
				utils.EmitKeyValue(word, chunk.FileName, w.PortNumber, w.NumReducers)
			}

		default:
			log.Printf("Unknown task description: %s", w.TaskDesc)
			return fmt.Errorf("unknown task description: %s", w.TaskDesc)
		}
	}

	// Send a response back to the client
	return stream.SendAndClose(&pb.FileInfo{Location: fmt.Sprintf("mapResults/%d", w.PortNumber)})
}

func (w *Worker) SendToReducer(ctx context.Context, req *pb.FileInfo) (*emptypb.Empty, error) {
	// The folder path is what we get in the req
	// The reducer will read all the files in the folder and then
	// reduce them by storing it in a map one by one
	// The map will then be used to write the final output to a file
	filePath := fmt.Sprintf("%s/%d.out", req.Location, w.PortNumber - 6000)
	log.Printf("Received file path: %s", filePath)

	// Lets read this file line by line, the key values are tab separated
	// We will split the line by tab and then store the key value pair in 
	// the worker map
	w.mu.Lock()
	err := utils.ReadIntermediateFile(filePath, w.ReducerList)
	w.mu.Unlock()
	
	if err != nil {
		log.Printf("Error reading intermediate file: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (w *Worker) Vomit(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	// // Reduce the intermediate data
	// w.mu.Lock()
	// output := w.ReducerList
	// w.mu.Unlock()

	// // Sort and reduce the data
	// sortedOutput := utils.SortKV(output)
	// reducedOutput := utils.ReduceByKey(sortedOutput, w.TaskDesc)

	// // Write the output to the file
	// outputPath := fmt.Sprintf("reducerResults/%s.out", w.TaskDesc)
	// err := utils.WriteOutputToFile(outputPath, reducedOutput)
	// if err != nil {
	// 	log.Printf("Error writing output to file: %v", err)
	// 	return nil, err
	// }

	return &emptypb.Empty{}, nil

}