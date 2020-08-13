package rpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/peer"

	"dcss/common/db"
	"dcss/core/entity"
	"dcss/core/proto"
)

var (
	_ proto.HeartbeatServer = &Heartbeat{}
)

type Heartbeat struct {
}

func New() *Heartbeat {
	return &Heartbeat{}
}

// worker向scheduler注册
func (h *Heartbeat) Registry(ctx context.Context, req *proto.RegistryReq) (*proto.Empty, error) {
	// 通过peer更新ip
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("peer is null, registry failed")
	}
	ip, _, _ := net.SplitHostPort(p.Addr.String())
	req.Ip = ip
	addr := fmt.Sprintf("%s:%d", req.Ip, req.Port)

	// 校验此host是否存在，若存在，则更新
	host := &entity.Host{}
	err := host.GetByAddr(addr, db.Db())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// todo create host
		}
		// todo print log
	} else {
		// todo update host
	}

	// 将此host绑定到对应的主机组，若host_group不为空的时候
	if req.HostGroup != "" {

	}
	return &proto.Empty{}, nil
}

func (h *Heartbeat) Heartbeat(ctx context.Context, req *proto.HeartbeatReq) (*proto.Empty, error) {
	// 通过peer更新ip
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("peer is null, update heartbeat failed")
	}
	ip, _, _ := net.SplitHostPort(p.Addr.String())
	addr := fmt.Sprintf("%s:%d", ip, req.Port)

	// 更新心跳时间和正在运行的任务
	host := &entity.Host{
		Addr:               addr,
		RunningTasks:       strings.Join(req.GetRunningTask(), ","),
		LastUpdateTimeUnix: time.Now().Unix(),
	}
	err := host.UpdateHeartbeatByAddr(db.Db())
	if err != nil {
		// todo print log
		return nil, errors.New("update heartbeat failed")
	}
	return &proto.Empty{}, nil
}
