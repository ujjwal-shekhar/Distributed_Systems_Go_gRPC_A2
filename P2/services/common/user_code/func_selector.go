package usercode

import (
	"github.com/ujjwal-shekhar/mapreduce/services/common/user_code/Task1"
	"github.com/ujjwal-shekhar/mapreduce/services/common/user_code/template"
)

// UserTaskDetails struct to hold task-specific details
type UserTaskDetails struct {
	Mapper                func([]common.KV, chan common.KV)
	Reducer               func([]common.KV, chan common.KV)
	KV_inType             common.KV
	KV_intermediateType   common.KV
	KV_outType            common.KV
}

// GetTaskDetails returns task-specific details based on the task name
func GetTaskDetails(taskName string) *UserTaskDetails {
	switch taskName {
	case "wordcount":
		return &UserTaskDetails{
			Mapper:                task1.Map,
			Reducer:               task1.Reduce,
			KV_inType:             &task1.KV_in{},
			KV_intermediateType:   &task1.KV_intermediate{},
			KV_outType:            &task1.KV_out{},
		}
	default:
		return nil
	}
}