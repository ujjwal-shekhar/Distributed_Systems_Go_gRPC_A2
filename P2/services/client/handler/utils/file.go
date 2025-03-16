package utils

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/ujjwal-shekhar/mapreduce/services/client/handler/client"
	pb "github.com/ujjwal-shekhar/mapreduce/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/mapreduce/services/common/utils"
)

// FileToMappersPipeline reads files from folderPath and sends chunks to chunkChannel.
func FileToMappersPipeline(folderPath string, chunkChannel chan ChunkMetadata) error {
	files, err := os.ReadDir(helpers.DATASET_PATH)
	if err != nil {
		return fmt.Errorf("error reading folder %s: %w", folderPath, err)
	}

	log.Println("Files in folder: ", files, len(files))

	// Process each file in a separate goroutine
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		go ReadFileInChunks(folderPath+"/"+file.Name(), chunkChannel)
		log.Println("File sent to chunk channel: ", file.Name())
	}

	return nil
}

func ReadFileInChunks(filePath string, chunkChannel chan ChunkMetadata) {
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

	for {
		n, err := reader.Read(buffer) // Read CHUNK_SIZE bytes at a time
		log.Println("Read chunk of size: ", n)
		// log.Printf("Read chunk: %s\n", buffer[:n])
		if n > 0 {
			chunkData := make([]byte, n)
			copy(chunkData, buffer[:n])
			
			chunkChannel <- ChunkMetadata{
				DocumentID: documentID,
				ChunkID:    chunkID,
				ChunkData:  chunkData,
			}
			
			log.Printf("Sent Chunk: %d, Size: %d bytes to the channel\n", chunkID, n)
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
}

func SendChunksToMappers(chunkChannel chan ChunkMetadata, master *client.Master) {
	mapperCount := len(master.MapperServers)
	if mapperCount == 0 {
		log.Println("No mapper servers available. Exiting chunk distribution.")
		return
	}
	
	i := 0 // Round-robin index
	var mu sync.Mutex
	for chunk := range chunkChannel {
		mu.Lock()
		mapper := master.MapperServers[i]

		// Send chunk to the mapper using gRPC
		_, err := mapper.FileUpload(context.Background(), &pb.FileChunk{
			Chunk:  		chunk.ChunkData,
			ChunkNumber:    int32(chunk.ChunkID),
			FileName:		chunk.DocumentID,
		})

		if err != nil {
			log.Printf("Failed to send chunk %d of document %s to mapper %d: %v", 
				chunk.ChunkID,  
				strings.Split(chunk.DocumentID, "/")[len(strings.Split(chunk.DocumentID, "/"))-1],
				i, err)
		} else {
			log.Printf("Successfully sent chunk %d of document %s to mapper %d", 
			chunk.ChunkID,  
			strings.Split(chunk.DocumentID, "/")[len(strings.Split(chunk.DocumentID, "/"))-1],
			i)
		}

		// Round-robin selection of next mapper
		i = (i + 1) % mapperCount

		mu.Unlock()
	}
}