package model

type Task struct {
	TaskType      int      `json:"task_type" binding:"required"`                 // 任务类型
	TaskData      string   `json:"task_data" binding:"required"`                 // 任务数据
	Run           bool     `json:"run" `                                         // 是否可以自动调度  如果为false则只能手动或者被其他任务依赖运行
	Creator       string   `json:"creator"`                                      // 创建人
	HostGroupID   string   `json:"host_group_id" binding:"required,len=18"`      // 主机组ID
	CronExpr      string   `json:"cron_expr" binding:"required,max=1000"`        // 执行任务表达式
	Timeout       int      `json:"timeout" binding:"required,min=-1"`            // 任务超时时间 (s) -1 no limit
	AlarmUserIds  []string `json:"alarm_user_ids" binding:"required,max=10"`     // 报警用户 最多十个多个用户
	RoutePolicy   int      `json:"route_policy" binding:"required,min=1,max=4"`  // how to select a run worker from host_group
	ExpectCode    int      `json:"expect_code"`                                  // expect task return code. if not set 0 or 200
	ExpectContent string   `json:"expect_content"`                               // expect task return content. if not set do not check
	AlarmStatus   int      `json:"alarm_status" binding:"required,min=-2,max=1"` // alarm when task run success or fail or all all:-2 failed: -1 success: 1
	Remark        string   `json:"remark" binding:"max=100"`
}

type CreateTask struct {
	Name string `json:"name"`
	Task
}
