package scheduler

import (
	"encoding/json"

	"dcss/common/db"
	"dcss/core/entity"
)

type EventType uint8

const (
	AddEvent EventType = iota + 1
	ChangeEvent
	DeleteEvent
	RunEvent
	KillEvent
)

const (
	TaskEventChannel = "task.event"
)

type Event struct {
	TaskID int
	Type   EventType
}

func (s *Scheduler) PushEvent(event []byte) {
	s.redis.Publish(TaskEventChannel, event)
}

func (s *Scheduler) recvEvent() {
	pubSub := s.redis.Subscribe(TaskEventChannel)
	for {
		select {
		case msg := <-pubSub.Channel():
			go s.dealEvent([]byte(msg.Payload))
		case <-s.close:
			break
		}
	}
}

func (s *Scheduler) dealEvent(data []byte) {
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		// todo print log
		return
	}
	switch event.Type {
	case AddEvent:
		// 校验是否已经在内存中
		// 若存在内存中，则直接return
		if s.isTaskExist(event.TaskID) {
			return
		}

		// 加载此任务
		task := &entity.Task{
			ID: event.TaskID,
		}
		err := task.GetByID(db.Db())
		if err != nil {
			return
		}

		// 转换task为运行时task
		// 将task加入到缓存当中
		// 根据是否自动调度决定task是否启动
		s.addTask(s.generateTask(task))
	}
}
