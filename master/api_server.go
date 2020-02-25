package master

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  api_server
 * @Version: 1.0.0
 * @Date: 2020/2/24 4:29 下午
 */

var (
	G_apiServer *apiServer
)

type apiServer struct {
	httpServer *http.Server
}

func InitApiServer() (err error) {
	var (
		mux      *http.ServeMux
		listener net.Listener
	)
	//注册路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	//启动http服务器
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}
	httpServer := &http.Server{
		Handler:      mux,
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
	}
	G_apiServer = &apiServer{httpServer: httpServer}
	go G_apiServer.httpServer.Serve(listener)
	return
}

func handleJobSave(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		//todo, print err
		return
	}
}
