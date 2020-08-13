package scheduler

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorhill/cronexpr"

	"dcss/common/db"
	"dcss/core/entity"
)

type Scheduler struct {
	sync.RWMutex
	redis *redis.Client
	tasks map[int]*task
	close chan struct{}
}

func New(redis *redis.Client) *Scheduler {
	s := &Scheduler{
		redis: redis,
		tasks: make(map[int]*task),
	}
	s.init()
	return s
}

func (s *Scheduler) init() {
	// 加载所有任务
	tasks, err := (&entity.Task{}).DescribeTasks(db.Db())
	if err != nil {
		// todo
		panic(err)
	}

	// 转换task并加入到cache当中
	for _, task := range tasks {
		s.addTask(s.generateTask(task))
	}

	// 启动频道订阅
	go s.recvEvent()
}

func (s *Scheduler) addTask(task *task) {
	// 存在性检测
	// 若存在且为自动调度，则需要先关闭任务循环
	task, ok := s.getTask(task.id)
	if ok && task.isAuto() {
		close(task.close)
		delete(s.tasks, task.id)
	}

	// task加入调度队列
	s.Lock()
	s.tasks[task.id] = task
	s.Unlock()
	if task.isAuto() {
		go s.scheduleTask(task.id)
	}
}

// 核心任务调度函数
func (s *Scheduler) scheduleTask(taskId int) {
	// 获取任务、解析表达式
	task, exist := s.getTask(taskId)
	if !exist {
		return
	}
	expr, err := cronexpr.Parse(task.cronExpr)
	if err != nil {
		return
	}

	// 设置上一次和下一次时间
	var (
		// 上一次开始执行时间
		last time.Time
		// 下一次执行时间
		next time.Time
	)
	last = time.Now()

	// 计算出锁的续约时间
	task.cronSub = expr.Next(last).Sub(last) / 4
	if task.cronSub > time.Second*30 {
		task.cronSub = time.Second * 30
	}
	for {
		next = expr.Next(last)
		select {
		case <-task.close:
			return
		case <-time.After(next.Sub(last)):
			last = next
			if task.canRun {
				go task.run()
			}
		}
	}
}

// task是否已经存在于内存中
func (s *Scheduler) isTaskExist(taskId int) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.tasks[taskId]
	return ok
}

// 获取给定id的task
func (s *Scheduler) getTask(taskId int) (*task, bool) {
	s.RLock()
	defer s.RUnlock()
	task, ok := s.tasks[taskId]
	return task, ok
}

// 删除一个task
func (s *Scheduler) deleteTask(taskId int) {
	task, exist := s.getTask(taskId)
	if exist {
		s.Lock()
		delete(s.tasks, taskId)
		s.Unlock()
		if task.close != nil {
			close(task.close)
		}
	}
}

// 给定entity的task，生成运行时的task
func (s *Scheduler) generateTask(t *entity.Task) *task {
	return &task{
		id:       t.ID,
		name:     t.Name,
		cronExpr: t.CronExpr,
		taskType: t.TaskType,
		close:    make(chan struct{}),
		canRun:   t.Run,
		redis:    s.redis,
		once:     sync.Once{},
	}
}
