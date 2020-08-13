package http

import "dcss/app/scheduler"

type Engine struct {
	scheduler *scheduler.Scheduler
}

func New(scheduler *scheduler.Scheduler) *Engine {
	return &Engine{scheduler: scheduler}
}
