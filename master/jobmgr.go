package master

import (
	"encoding/json"
	"time"

	"github.com/samuel/go-zookeeper/zk"

	"crontab/common"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  jobmgr
 * @Version: 1.0.0
 * @Date: 2020/2/24 5:26 下午
 */

var (
	G_jobMgr *jobMgr
)

type jobMgr struct {
	conn    *zk.Conn
	jobPath string
}

func InitJobMgr() error {
	conn, _, err := zk.Connect(G_config.ZkIps, time.Duration(G_config.ZkTimeout)*time.Millisecond)
	if err != nil {
		return err
	}
	G_jobMgr = &jobMgr{
		conn:    conn,
		jobPath: G_config.JobPath,
	}
	return nil
}

func (m *jobMgr) SaveJob(job *common.Job) error {
	var (
		jobKey   string
		jobValue []byte
		err      error
	)
	jobKey = m.jobPath + job.Name
	if jobValue, err = json.Marshal(job); err != nil {
		return err
	}
	if _, _, err = m.conn.Get(jobKey); err != nil {
		if err == zk.ErrNoNode {
			_, err = m.conn.Create(jobKey, jobValue, 0, zk.WorldACL(31))
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if _, err = m.conn.Set(jobKey, jobValue, 0); err != nil {
		return err
	}
	return nil
}
