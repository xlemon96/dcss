package main

import (
	"flag"
	"fmt"
	"runtime"

	"crontab/master"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2020/2/24 4:27 下午
 */

var (
	filename string
)

func main() {
	var (
		err error
	)
	//解析命令行参数
	initArgs()
	//设置线程数目
	initEnv()
	//加载配置
	if err = master.InitConfig(filename); err != nil {
		goto ERR
	}
	//初始化zk客户端
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}
	//初始化http服务器
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}
	//正常退出
	return
//异常处理
ERR:
	fmt.Println(err)
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initArgs() {
	flag.StringVar(&filename, "config", "./master.json", "指定配置文件")
	flag.Parse()
}