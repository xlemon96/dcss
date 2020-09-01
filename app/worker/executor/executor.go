package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"dcss/core/proto"
	"dcss/core/running"
)

const (
	DefaultExitCode int = -1
)

type TaskExecutor interface {
	Run(ctx context.Context) (out io.ReadCloser)
	Type() string
}

func GetTaskExecutor(t *proto.TaskReq) (TaskExecutor, error) {
	switch running.TaskType(t.TaskType) {
	case running.Code:
		var code Code
		err := json.Unmarshal(t.TaskData, &code)
		if err != nil {
			return nil, err
		}
		code.LangDesc = code.Lang.String()
		return code, err
	case running.API:
		var api API
		err := json.Unmarshal(t.TaskData, &api)
		if err != nil {
			return nil, err
		}
		if api.Header == nil {
			api.Header = make(map[string]string)
		}
		return api, err
	default:
		err := fmt.Errorf("unsupport task type %d", t.TaskType)
		return nil, err
	}
}
