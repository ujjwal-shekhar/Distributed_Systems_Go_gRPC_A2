package handler

import (
	"context"
	"log"
	"strconv"
	"sync"

	emptypb "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/ujjwal-shekhar/mapreduce/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/mapreduce/services/common/user_code"
	"github.com/ujjwal-shekhar/mapreduce/services/common/user_code/template"

	"github.com/ujjwal-shekhar/mapreduce/services/common/user_code/Task1"
	// "github.com/ujjwal-shekhar/mapreduce/services/common/user_code/Task2"
)

type Worker struct {
	pb.UnimplementedFileTransferServer
	
	// Mu 				*sync.Mutex
	// Reducers 		[]pb.FileTransferClient
	IsMapper		bool
	TaskDesc		string
	TaskChannel 	chan common.KV
	WorkFunc 		func([]common.KV, chan common.KV)
	UserDetails     *usercode.UserTaskDetails
	NumReducers		int
}

func NewMapper(taskDetails *usercode.UserTaskDetails, taskType string, numReducers int) *Worker {
	return &Worker {
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
		IsMapper: false,
		TaskDesc: taskType,
		TaskChannel: make(chan common.KV),
		WorkFunc: taskDetails.Reducer,
		UserDetails: taskDetails,
		NumReducers: numReducers,
	}
}

func (w *Worker) FileUpload(ctx context.Context, in *pb.FileChunk) (*emptypb.Empty, error) {
	log.Printf("Received chunk %v %v", in.FileName, in.ChunkNumber)

	// The client would have made this call to the mapper
	// So, we need to create a key-value pair of the input type
	// and push it to the task channel so the mapper function
	// emits the intermediate key-value pairs to another outchannel
	kv := w.UserDetails.KV_inType

	if w.TaskDesc == "wordcount"{
		kv.SetKeyVal(
			task1.Key_in{Item: in.FileName},
			task1.Value_in{Item: string(in.Chunk)},
		)
	} else if w.TaskDesc == "invertedindex"{
		// kv.SetKeyVal(
		// 	task2.Key_in{Item: in.FileName},
		// 	task2.Value_in{Item: string(in.Chunk)},
		// )
	} else {
		log.Fatalf("Failed to create key-value pair")
	}
	w.TaskChannel <- kv

	return &emptypb.Empty{}, nil
}

func (w *Worker) SendKV(ctx context.Context, in *pb.KVPair) (*emptypb.Empty, error) {
	log.Printf("Received key-value pair: %v", in)

	// A mapper would have made this call to the reducer
	// So, we will funnel all the intermediate key value pairs
	// to the task channel so reducers can process them
	// and then it will be thrown at the outchannel
	kv := w.UserDetails.KV_intermediateType

	if w.TaskDesc == "wordcount"{
		val, err := strconv.Atoi(string(in.Value))
		if err != nil {
			log.Fatalf("Failed to convert value to int: %v", err)
		}

		kv.SetKeyVal(
			task1.Key_intermediate{Item: in.Key},
			task1.Value_intermediate{Item: int32(val)},
		)
	}  else if w.TaskDesc == "invertedindex"{
		// val, err := strconv.Atoi(string(in.Chunk))
		// if err != nil {
		// 	log.Fatalf("Failed to convert value to int: %v", err)
		// }
		// kv.SetKeyVal(
		// 	task2.Key_intermediate{Item: in.FileName},
		// 	task2.Value_intermediate{Item: in.StringWorkLoad},
		// )
	} else {
		log.Fatalf("Failed to create key-value pair")
	}
	w.TaskChannel <- kv

	return &emptypb.Empty{}, nil
}