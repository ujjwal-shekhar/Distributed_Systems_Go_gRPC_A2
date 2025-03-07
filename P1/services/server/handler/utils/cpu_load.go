package utils

import (
	"math/rand"
	"runtime"
)

func FakeTask(iters int32) {
	var result float64
	for i := int32(0); i < iters; i++ {
		result += rand.Float64() * rand.Float64()
		if result > 1000 {
			result = 0
		}
	}
	runtime.KeepAlive(result) // Prevents optimization
}
