package task1

import (
	"strings"

	"github.com/ujjwal-shekhar/mapreduce/services/common/user_code/template"
)

// Map function for the "wordcount" task
func Map(kvs []common.KV, resultChan chan common.KV) {
	for _, kv := range kvs {
		// Type assertion to convert common.KV to KV_in
		// bro this took too much time to understand
		kvIn := kv.(*KV_in)
		words := strings.Fields(kvIn.Value.Item)

		// Emit intermediate key-value pairs
		for _, word := range words {
			resultChan <- &KV_intermediate{
				Key:   Key_intermediate{Item: word},
				Value: Value_intermediate{Item: 1},
			}
		}
	}
}
