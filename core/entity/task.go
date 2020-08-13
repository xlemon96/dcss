package entity

import "github.com/jinzhu/gorm"

type Task struct {
	ID            int    // 主键id
	Name          string // 任务名称
	TaskType      int    // 任务类型
	TaskData      string // 任务执行需要的数据
	Run           bool   // 是否自动调度，若为false，则手动执行
	Creator       string // 创建者
	HostGroupID   string // 所属主机组id
	CronExpr      string // 定时任务表达式
	Timeout       int    // 定时任务超时时间
	AlarmUserIds  string // 报警用户
	RoutePolicy   int    // 路由策略 1:Random 2:RoundRobin 3:Weight 4:LeastTask
	ExpectCode    int    // 期望状态码 CODE默认为0 HTTP默认为200
	ExpectContent string // 期望输出内容的文本
	AlarmStatus   int    // 报警策略 1:任务运行结束 2:任务运行失败 3:任务运行成功
	Remark        string // 备注
	CreateTime    string // 创建时间
	UpdateTime    string // 更新时间
}

func (t *Task) CreateTask(db *gorm.DB) (int, error) {
	return t.ID, db.Create(t).Error
}

func (t *Task) DescribeTasks(db *gorm.DB) ([]*Task, error) {
	tasks := make([]*Task, 0)
	err := db.Find(&tasks).Error
	return tasks, err
}

func (t *Task) DescribeTaskByID(db *gorm.DB) error {
	db = db.Where("id = ?", t.ID)
	return db.Find(t).Error
}
