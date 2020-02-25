package common

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  protocol
 * @Version: 1.0.0
 * @Date: 2020/2/24 5:23 下午
 */

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}
