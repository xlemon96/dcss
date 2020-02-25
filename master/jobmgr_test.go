package master

import (
	"fmt"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"

	"crontab/common"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  jobmgr_test
 * @Version: 1.0.0
 * @Date: 2020/2/24 6:20 下午
 */

var jobmgr *jobMgr

func init() {
	conn, _, err := zk.Connect([]string{"localhost"}, 5000*time.Second)
	if err != nil {
		return
	}
	jobmgr = &jobMgr{
		conn:    conn,
		jobPath: "/cron/jobs/",
	}
}

func TestJobMgr_SaveJob(t *testing.T) {
	job := &common.Job{
		Name:     "test",
		Command:  "test",
		CronExpr: "test",
	}
	err := jobmgr.SaveJob(job)
	if err != nil {
		fmt.Println(err)
	}
}
