package handler

import (
	"log"
	"io"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/ujjwal-shekhar/mapreduce/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/mapreduce/services/server/handler/utils"
)

type Worker struct {
	*pb.UnimplementedFileTransferServer
	
	IsMapper		bool
	Chunks 			[]utils.ChunkMetadatad
}

func (w *Worker) FileUpload(stream pb.FileTransfer_FileUploadServer) error {
	for {
		req, err := stream.Recv() // Receive a chunk from the client
		if err == io.EOF {
			// Done receiving chunks
			log.Println("File upload complete.")
			return stream.SendAndClose(&emptypb.Empty{}) // Acknowledge completion
		}
		if err != nil {
			log.Printf("Error receiving file chunk: %v", err)
			return err
		}

		log.Printf("Received chunk %d of file %s, size: %d bytes\n", req.ChunkNumber, req.FileName, len(req.Chunk))

		// Store the chunk
		utils.StoreChunk(req.FileName, req.ChunkNumber, req.Chunk)
	}
}

func NewMapper() *Worker {
	return &Worker{
		IsMapper: true,
		Chunks:   make([]utils.ChunkMetadata, 0),
	}
}

func NewReducer() *Worker {
	return &Worker{
		IsMapper: false,
		Chunks:   make([]utils.ChunkMetadata, 0),
	}
}

