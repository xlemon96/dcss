package http

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorhill/cronexpr"

	"dcss/app/scheduler"
	"dcss/app/scheduler/http/model"
	"dcss/common/db"
	"dcss/core/entity"
	"dcss/core/running"
	"dcss/util"
)

func (e *Engine) CreateTask(c *gin.Context) {
	// 解析参数
	var req model.CreateTask
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(running.ErrBadRequest, nil)
		return
	}

	// 校验定时表达式
	_, err = cronexpr.Parse(req.CronExpr)
	if err != nil {
		c.JSON(running.ErrCronExpr, err)
		return
	}

	// 任务写入数据库
	task := &entity.Task{
		Name:          req.Name,
		TaskType:      req.TaskType,
		TaskData:      req.TaskData,
		Run:           req.Run,
		Creator:       req.Creator,
		HostGroupID:   req.HostGroupID,
		CronExpr:      req.CronExpr,
		Timeout:       req.Timeout,
		AlarmUserIds:  strings.Join(req.AlarmUserIds, ";"),
		RoutePolicy:   req.RoutePolicy,
		ExpectCode:    req.ExpectCode,
		ExpectContent: req.ExpectContent,
		AlarmStatus:   req.AlarmStatus,
		Remark:        req.Remark,
		CreateTime:    util.TimeToString(time.Now()),
		UpdateTime:    util.TimeToString(time.Now()),
	}
	taskId, err := task.CreateTask(db.Db())
	if err != nil {
		c.JSON(running.ErrInternalServer, err)
		return
	}

	// 向redis推送消息
	event := scheduler.Event{
		TaskID: taskId,
		Type:   scheduler.AddEvent,
	}
	bEvent, err := json.Marshal(event)
	if err != nil {
		c.JSON(running.ErrInternalServer, err)
		return
	}
	e.scheduler.PushEvent(bEvent)
	c.JSON(running.Success, nil)
}
