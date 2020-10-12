package rpc

import (
	"context"
	"fmt"
	"io"

	"dcss/app/worker/executor"
	"dcss/core/proto"
)

var (
	_ proto.TaskServer = &Task{}
)

type Task struct {
}

func New() *Task {
	return &Task{}
}

func (t *Task) RunTask(req *proto.TaskReq, stream proto.Task_RunTaskServer) error {
	e, err := executor.GetTaskExecutor(req)
	if err != nil {
		err = stream.Send(&proto.TaskResp{Resp: []byte(err.Error())})
		if err != nil {
			// todo print log
		}
		return nil
	}
	ctx, _ := context.WithCancel(stream.Context())
	out := e.Run(ctx)
	// 缓存任务，退出时删除任务
	defer out.Close()
	// 循环读取脚本执行的输出
	var buf = make([]byte, 1024)
	for {
		n, err := out.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			err = stream.Send(&proto.TaskResp{Resp: []byte(err.Error() +
				fmt.Sprintf("%3d", executor.DefaultExitCode))})
			if err != nil {
			}
			return nil
		}
		if n > 0 {
			resp := proto.TaskResp{Resp: buf[:n]}
			err = stream.Send(&resp)
			if err != nil {
				return nil
			}
		}
	}
}
