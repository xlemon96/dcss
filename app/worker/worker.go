package worker

import (
	"context"
	"sync"
)

type Task struct {
	sync.RWMutex
	running map[string]context.CancelFunc
}
