package entity

import "dcss/core/running"

type Log struct {
	Name           string               `json:"name"`                 // task log
	RunByTaskID    string               `json:"runby_taskid"`         // run taskid
	StartTime      int64                `json:"start_time"`           // ms
	StartTimeStr   string               `json:"start_timestr"`        //
	EndTime        int64                `json:"end_time"`             // ms
	EndTimeStr     string               `json:"end_timestr"`          //
	TotalRunTime   int                  `json:"total_runtime"`        // ms
	Status         int                  `json:"status"`               // 任务运行结果 -1 失败 1 成功
	TaskResps      []*running.TaskResp  `json:"task_resps,omitempty"` // 任务执行过程日志
	Trigger        running.Trigger      `json:"trigger"`              // 任务触发
	Triggerstr     string               `json:"trigger_str"`          // 任务触发
	ErrCode        int                  `json:"err_code"`             // err code
	ErrMsg         string               `json:"err_msg"`              // 错误原因
	ErrTasktype    running.TaskRespType `json:"err_tasktype"`         // err task type
	ErrTaskTypeStr string               `json:"err_tasktypestr"`      // 1 主任务 2 父任务 3 子任务
	ErrTaskID      string               `json:"err_taskid"`           // task failed id
	ErrTask        string               `json:"err_task"`             // task failed id
}
