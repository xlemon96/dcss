package running

// 任务类型
type TaskType uint8

const (
	Code TaskType = iota + 1
	API
)

func (t TaskType) String() string {
	switch t {
	case Code:
		return "code"
	case API:
		return "api"
	default:
		return "unknown"
	}
}
