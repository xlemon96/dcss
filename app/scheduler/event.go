package scheduler

import "encoding/json"

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

func (s *Scheduler) recvEvent() {
	pubSub := s.redis.Subscribe(TaskEventChannel)
	for msg := range pubSub.Channel() {
		go s.dealEvent([]byte(msg.Payload))
	}
}

func (s *Scheduler) dealEvent(data []byte) {
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		// todo print log
		return
	}
}
