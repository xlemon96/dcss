package scheduler

import (
	"sync"

	"github.com/go-redis/redis"
)

type Scheduler struct {
	sync.RWMutex
	redis *redis.Client
	tasks map[string]*task
}

func New(redis *redis.Client) *Scheduler {
	s := &Scheduler{
		redis: redis,
		tasks: make(map[string]*task),
	}
	s.init()
	return s
}

func (s *Scheduler) init() {
	// 加载所有任务

	// 启动频道订阅
}

func (s *Scheduler) addTask() {

}

// task是否已经存在于内存中
func (s *Scheduler) isTaskExist(taskId string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.tasks[taskId]
	return ok
}
