package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func SplitChunksToWords(chunk []byte) []string {
	// Split the chunk into words
	words := strings.Fields(string(chunk))
	return words
}

// emitKeyValue writes a key-value pair to the appropriate intermediate file.
func EmitKeyValue(key, value string, portNumber, reduceTasks int) {
	hashKey := hash(key) % reduceTasks
	outputFile := fmt.Sprintf("mapResults/%d/%d.out", portNumber, hashKey)

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(outputFile), os.ModePerm); err != nil {
		log.Printf("Error creating directory: %v", err)
		return
	}

	// Append the key-value pair to the intermediate file
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file %s: %v", outputFile, err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%s\t%s\n", key, value)); err != nil {
		log.Printf("Error writing to file %s: %v", outputFile, err)
		return
	}
}

type KV struct {
	Key   string
	Value string
}

func ReadIntermediateFile(filePath string, kv []KV) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		return err
	}
	defer file.Close()

	// Use bufio for proper buffered reading
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("Finished reading file:", filePath)
			} else {
				log.Printf("Error reading file %s: %v", filePath, err)
				return err
			}
			break
		}

		// Split the line into key and value
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			log.Printf("Invalid line format: %s", line)
			continue
		}
		kv = append(kv, KV{Key: parts[0], Value: parts[1]})
	}

	return nil
}

type ReducedKV struct {
	Key   string
	Value []string
}

func WriteOutputToFile(filePath string, kv []ReducedKV) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file %s: %v", filePath, err)
		return err
	}
	defer file.Close()

	// Write the key-value pairs to the output file
	for _, rkv := range kv {
		if _, err := file.WriteString(fmt.Sprintf("%s : %s\n", rkv.Key, strings.Join(rkv.Value, ","))); err != nil {
			log.Printf("Error writing to file %s: %v", filePath, err)
			return err
		}
	}

	return nil
}
