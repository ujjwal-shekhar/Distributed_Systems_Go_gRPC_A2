package utils

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ujjwal-shekhar/mapreduce/services/client/handler/client"
	pb "github.com/ujjwal-shekhar/mapreduce/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/mapreduce/services/common/utils"
)

// FileToMappersPipeline reads files from folderPath and sends chunks to chunkChannel.
func FileToMappersPipeline(folderPath string, master *client.Master, wg *sync.WaitGroup) error {
	files, err := os.ReadDir(helpers.DATASET_PATH)
	if err != nil {
		return fmt.Errorf("error reading folder %s: %w", folderPath, err)
	}

	log.Println("Files in folder: ", files, len(files))

	// Process each file in a separate goroutine
	for i, file := range files {
		if file.IsDir() {
			continue
		}

		// wg.Add(1)
		// go ReadFileInChunks(folderPath+"/"+file.Name(), i, master, wg)
		wg.Add(1)
		go func(fileName string, mapperId int) {
			defer wg.Done()
			ReadFileInChunks(folderPath+"/"+fileName, mapperId, master)
		}(file.Name(), i)
	
		log.Println("File sent to mapper handler: ", file.Name())
	}

	return nil
}

func ReadFileInChunks(filePath string, mapperId int, master *client.Master) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		return
	}
	defer file.Close()

	// Use bufio for proper buffered reading
	reader := bufio.NewReader(file)
	buffer := make([]byte, helpers.CHUNK_SIZE)
	chunkID := 0
	documentID := filePath

	// Initiate the stream
	stream, err := master.MapperServers[mapperId].SendToMapper(context.Background())
	if err != nil {
		log.Printf("Error sending to mapper %d: %v", mapperId, err)
		return
	}

	for {
		n, err := reader.Read(buffer) ; if n > 0 {
			chunkData := make([]byte, n)
			copy(chunkData, buffer[:n])
			
			// Send chunk to the mapper
			chunk := &pb.FileChunk{
				Chunk:       chunkData,
				ChunkNumber: int32(chunkID),
				FileName:    documentID,
			}

			if err := stream.Send(chunk); err != nil {
				log.Printf("Error sending chunk to mapper %d: %v", mapperId, err)
				return
			}
			
			log.Printf("Sent Chunk: %d, Size: %d bytes to the mapper %d\n", chunkID, n, mapperId)
			chunkID++
		}

		if err != nil {
			if err.Error() == "EOF" {
				log.Println("Finished reading file:", documentID)
			} else {
				log.Printf("Error reading file %s: %v", filePath, err)
			}
			break
		}
	}

	// Close the stream and get intermediate info
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("Error closing stream to mapper %d: %v", mapperId, err)
		return
	}

	log.Printf("Intermediate info from mapper %d: %s\n", mapperId, reply.Location)

	// Send this to the reducer
	for i := 0; i < len(master.ReducerServers); i++ {
		reducer := master.ReducerServers[i]
		_, err := reducer.SendToReducer(context.Background(), &pb.FileInfo{
			Location: reply.Location,
		})

		if err != nil {
			log.Printf("Error sending to reducer %d: %v", i, err)
			return
		}
	}
}

func CountFilesInFolder(folderPath string) int {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("Error reading folder %s: %v", folderPath, err)
	}

	count := 0
	for _, file := range files {
		if !file.IsDir() {
			count++
		}
	}

	return count
}