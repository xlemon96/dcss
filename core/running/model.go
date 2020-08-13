package running

// Trigger return how to trigger run task
type Trigger uint8

const (
	// Auto cron run task
	Auto Trigger = iota + 1
	// Manual trigger run task
	Manual
)

func (t Trigger) String() string {
	switch t {
	case Auto:
		return "自动触发"
	case Manual:
		return "手动触发"
	default:
		return "UnKnown"
	}
}

type TaskRespType uint8

const (
	// MasterTask task as master run
	MasterTask TaskRespType = iota + 1
	// ParentTask task as a task's parent task run
	ParentTask
	// ChildTask task as a task's child task run
	ChildTask
)

// TaskResp run task resp message
type TaskResp struct {
	TaskID      string       `json:"task_id"`
	Task        string       `json:"task"`
	LogData     string       `json:"resp_data"`    // task run log data
	Code        int          `json:"code"`         // return code
	TaskType    TaskRespType `json:"task_type"`    // 1 主任务 2 父任务 3 子任务
	TaskTypeStr string       `json:"task_typestr"` // 1 主任务 2 父任务 3 子任务
	RunHost     string       `json:"run_host"`     // task run host
	Status      string       `json:"status"`       // task status finish,fail, cancel
}
