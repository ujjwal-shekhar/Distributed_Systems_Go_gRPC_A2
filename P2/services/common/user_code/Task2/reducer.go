package task2

import (
	"github.com/ujjwal-shekhar/mapreduce/services/common/user_code/template"
)

// Reduce function for the "invertedindex" task
func Reduce(intermediates []common.KV, resultChan chan common.KV) {
	// Grouping the intermediate key-value pairs by key
	grouped := make(map[string][]string)
	for _, kv := range intermediates {
		kvIntermediate := kv.(*KV_intermediate)
		grouped[kvIntermediate.Key.Item] = append(grouped[kvIntermediate.Key.Item], kvIntermediate.Value.Item)
	}

	// Emit final key-value pairs
	for key, values := range grouped {
		resultChan <- &KV_out{
			Key:   Key_out{Item: key},
			Value: Value_out{Item: values},
		}
	}
}