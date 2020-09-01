package conn

import (
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"

	"dcss/core/running"
)

const (
	defaultRPCTimeout = time.Second * 3
	defaultMaxRetry   = 3
)

var (
	manager *gRPCConnManager
)

func init() {
	manager = &gRPCConnManager{
		conn: make(map[string]*grpc.ClientConn),
	}
}

type gRPCConnManager struct {
	sync.RWMutex
	conn map[string]*grpc.ClientConn
}

func (g *gRPCConnManager) getConn(addr string) *grpc.ClientConn {
	g.Lock()
	conn, exist := g.conn[addr]
	g.Unlock()
	if exist && conn.GetState() == connectivity.Ready {
		return conn
	}
	if conn != nil {
		conn.Close()
	}
	delete(g.conn, addr)
	return nil
}

func (g *gRPCConnManager) addConn(addr string, conn *grpc.ClientConn) {
	g.Lock()
	g.conn[addr] = conn
	g.Unlock()
}

func GetRpcConn(next running.Next) (*grpc.ClientConn, error) {
	var (
		err  error
		conn *grpc.ClientConn
	)
	for i := 0; i < defaultMaxRetry; i++ {
		host := next()
		if host == nil {
			err = errors.New("next host is nil")
			continue
		}
		conn = getRpcConn(host.Addr)
		if conn == nil {
			continue
		}
	}
	return nil, err
}

func getRpcConn(addr string) *grpc.ClientConn {
	conn := manager.getConn(addr)
	if conn != nil {
		return conn
	}
	var err error
	options := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(16 * 1024 * 1024)), // 16M
		grpc.WithBlock(),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoff.Config{MaxDelay: time.Second * 2}, MinConnectTimeout: time.Second * 2}),
	}
	options = append(options, grpc.WithInsecure())
	ctx, cancel := context.WithTimeout(context.Background(), defaultRPCTimeout)
	defer cancel()
	conn, err = grpc.DialContext(ctx, addr, options...)
	if err != nil {
		return nil
	}
	if conn.GetState() == connectivity.Ready {
		manager.addConn(addr, conn)
		return conn
	} else {
		conn.Close()
	}
	return nil
}
