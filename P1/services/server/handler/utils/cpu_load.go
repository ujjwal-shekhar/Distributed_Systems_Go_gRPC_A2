package utils

import (
	"math/rand"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"time"
)

func FakeTask(iters int32) {
	var result float64
	endTime := time.Now().Add(time.Duration(iters) * time.Second)

	for time.Now().Before(endTime) {
		result += rand.Float64() * rand.Float64()
		if result > 1000 {
			result = 0
		}
	}

	runtime.KeepAlive(result) // Prevents optimization
}

func GetCPULoad() float32 {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0
	}

	if len(percentages) > 0 {
		return float32(percentages[0])
	}
	return 0
}