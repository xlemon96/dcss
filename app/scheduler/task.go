package scheduler

import (
	"context"
	"sync"
	"time"

	pt "dcss/core/task"
	"github.com/go-redis/redis"
)

// 运行时的task描述结构体
type task struct {
	sync.RWMutex                    // 锁
	id           string             // 任务id
	name         string             // 任务名称
	cronExpr     string             // 定时任务表达式
	cronSub      time.Duration      // cronexpt sub
	close        chan struct{}      // 任务停止chan
	cancel       context.CancelFunc // store cancelfunc could cancel all task by this cancel
	//next      Next               // it save a func Next by route policy

	redis *redis.Client // redis客户端
	once  sync.Once     // once

	// 任务失败相关描述信息
	errTaskID   string          // 失败的任务id
	errTask     string          // 失败的子任务
	errCode     int             // 失败的任务code
	errMsg      string          // 失败的任务信息
	errTaskType pt.TaskRespType // 失败的任务类型
}

func (t *task) run() {
	// 获取锁
}
