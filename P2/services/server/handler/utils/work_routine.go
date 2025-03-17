package utils

import (
	"hash/fnv"
	"log"
	"slices"
	"sort"
	"strconv"
)

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

func SortKV(kv []KV) []KV{
	// Sort the key-value pairs by key
	sort.Slice(kv, func(i, j int) bool {
		return kv[i].Key < kv[j].Key
	})

	return kv
}

func ReduceByKey(sorted_kv []KV, taskDesc string) []ReducedKV {
	// Reduce the key-value pairs by key
	var reduced_kv []ReducedKV
	var currentKey string
	var currentValue []string = make([]string, 0)
	for _, kv := range sorted_kv {
		if kv.Key == currentKey {
			if taskDesc == "wordcount" {
				val, err := strconv.Atoi(currentValue[0])
				if err != nil {
					log.Fatalf("Error converting value to int: %v", err)
				}
				currentValue[0] = strconv.Itoa(val + 1)
			} else if taskDesc == "invertedindex" {
				// currentValue = append(currentValue, kv.Value)
				// Append if not already present
				if !slices.Contains(currentValue, kv.Value) {
					currentValue = append(currentValue, kv.Value)
				}
			}
		} else {
			if currentKey != "" {
				reduced_kv = append(reduced_kv, ReducedKV{currentKey, currentValue})
			}
			currentKey = kv.Key
			if taskDesc == "wordcount" {
				currentValue = []string{"1"}
			} else if taskDesc == "invertedindex" {
				currentValue = []string{kv.Value}
			}
		}
	}
	if currentKey != "" {
		reduced_kv = append(reduced_kv, ReducedKV{currentKey, currentValue})
	}

	return reduced_kv
}