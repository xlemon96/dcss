package rpc

import "dcss/core/proto"

var (
	_ proto.TaskServer = &Executor{}
)

type Executor struct {
}

func New() *Executor {
	return &Executor{}
}

func (e *Executor) RunTask(req *proto.TaskReq, stream proto.Task_RunTaskServer) error {

}
