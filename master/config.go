package master

import (
	"encoding/json"
	"io/ioutil"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  config
 * @Version: 1.0.0
 * @Date: 2020/2/24 5:02 下午
 */

var (
	G_config *config
)

type config struct {
	ApiPort         int      `json:"apiPort"`
	ApiReadTimeout  int      `json:"apiReadTimeout"`
	ApiWriteTimeout int      `json:"apiWriteTimeout"`
	ZkIps           []string `json:"zkIps"`
	ZkTimeout       int      `json:"zkTimeout"`
	JobPath         string   `json:"jobPath"`
}

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    config
	)
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	G_config = &conf
	return
}
