package worker

import (
	"context"
	"sync"
)

type Task struct {
	sync.RWMutex
	running map[string]context.CancelFunc
}

func New() *Task {
	return &Task{
		running: make(map[string]context.CancelFunc),
	}
}

func (t *Task) Add(id string, cancel context.CancelFunc) {
	t.Lock()
	t.running[id] = cancel
	t.Unlock()
}

func (t *Task) Del(id string) {
	t.Lock()
	delete(t.running, id)
	t.Unlock()
}

func (t *Task) Cancel(id string) {
	t.RLock()
	cancel, ok := t.running[id]
	t.RUnlock()
	if !ok {
		return
	}
	cancel()
}

func (t *Task) GetRunningTasks() []string {
	var runningTasks []string
	t.RLock()
	for taskName := range t.running {
		runningTasks = append(runningTasks, taskName)
	}
	t.RUnlock()
	return runningTasks
}
