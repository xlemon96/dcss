package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/go-redis/redis"

	"dcss/common/db"
	"dcss/core/conn"
	"dcss/core/entity"
	"dcss/core/proto"
	"dcss/core/running"
)

const (
	taskRunningLockPrefix       = "task:run:lock:"
	defaultTaskRunningLockValue = 1
)

// 运行时的task描述结构体
type task struct {
	sync.RWMutex               // 锁
	id           int           // 任务id
	name         string        // 任务名称
	cronExpr     string        // 定时任务表达式
	taskType     int           // 任务类型
	cronSub      time.Duration // 刷新锁的时间间隔
	close        chan struct{} // 任务停止chan
	canRun       bool          // 是否可以调度
	next         running.Next  // 调度机器

	redis *redis.Client // redis客户端
	once  sync.Once     // once

	// 任务失败相关描述信息
	errCode int    // 失败的任务code
	errMsg  string // 失败的任务信息
}

func (t *task) run() {
	// 判断任务是否已经在运行
	if t.isOnLock() {
		return
	}

	// 抢锁
	if !t.lock() {
		return
	}
	defer t.unLock()

	// 定时续约，保证锁不被释放
	stopLease := make(chan struct{})
	go func() {
		ticker := time.NewTicker(t.cronSub * 3 / 4)
		defer ticker.Stop()
		for {
			select {
			case <-stopLease:
				return
			case <-ticker.C:
				if t.cronSub >= time.Second {
					t.redis.Expire(t.generateLockKey(), t.cronSub)
				} else {
					t.redis.PExpire(t.generateLockKey(), t.cronSub)
				}
			}
		}
	}()
	defer close(stopLease)

	// 获取任务需要的执行数据
	task := entity.Task{ID: t.id}
	err := task.DescribeTaskByID(db.Db())
	if err != nil {

	}
	taskData, err := json.Marshal(task.TaskData)
	if err != nil {

	}

	// 获取rpc连接和req
	c, err := conn.GetRpcConn(t.next)
	if err != nil {
		return
	}
	runTaskReq := &proto.TaskReq{
		TaskId:   int64(t.id),
		TaskType: int32(t.taskType),
		TaskData: taskData,
	}

	// 获取worker rpc客户端并调用RunTask方法
	var (
		ctx    context.Context
		cancel context.CancelFunc
		output []byte
	)
	if task.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(task.Timeout))
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()
	workerClient := proto.NewTaskClient(c)
	respStream, err := workerClient.RunTask(ctx, runTaskReq)
	if err != nil {

	}
	for {
		resp, err := respStream.Recv()
		// todo
		if err != nil {
			if err == io.EOF {
				err = nil
				// 获取返回码
			}
		}
		output = append(output, resp.GetResp()...)
	}
}

func (t *task) isAuto() bool {
	if running.Trigger(t.taskType) == running.Auto {
		return true
	}
	return false
}

func (t *task) lock() bool {
	success, err := t.redis.SetNX(t.generateLockKey(),
		defaultTaskRunningLockValue,
		t.cronSub).Result()
	if err != nil {
		// todo print log
		return false
	}
	return success
}

func (t *task) unLock() {
	script := redis.NewScript(`
		if redis.call("get",KEYS[1]) == ARGV[1] then
			return redis.call("del",KEYS[1])
		else
			return 0
		end
	`)
	_, err := script.Run(t.redis,
		[]string{t.generateLockKey()},
		defaultTaskRunningLockValue).Result()
	if err != nil {
		// todo print log
	}
}

// 判断是否任务锁已经存在
func (t *task) isOnLock() bool {
	run, err := t.redis.Exists(t.generateLockKey()).Result()
	if err != nil {
		return true
	}
	if run == 0 {
		return false
	}
	return true
}

// 生成任务运行时锁对应的key
func (t *task) generateLockKey() string {
	return fmt.Sprintf(taskRunningLockPrefix+"%d", t.id)
}
