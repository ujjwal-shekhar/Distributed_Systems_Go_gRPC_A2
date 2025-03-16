package task1

import (
	"log"

	"github.com/ujjwal-shekhar/mapreduce/services/common/user_code/template"
)

// Reduce function for the "wordcount" task
func Reduce(intermediates []common.KV, resultChan chan common.KV) {
	wordCounts := make(map[string]int32)
	for _, kv := range intermediates {
		// Type assertion to convert common.KV to KV_intermediate
		intermediateKV := kv.(*KV_intermediate)
		word := intermediateKV.Key.Item
		count := intermediateKV.Value.Item
		wordCounts[word] += count
	}

	// Emit final key-value pairs
	for word, count := range wordCounts {
		resultChan <- &KV_out{
			Key:   Key_out{Item: word},
			Value: Value_out{Item: count},
		}
	}
}