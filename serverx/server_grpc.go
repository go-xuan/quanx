package serverx

import (
	"context"
	"fmt"
	"net"

	"github.com/go-xuan/utilx/errorx"
	"google.golang.org/grpc"
)

// NewGrpcServer 创建grpc服务
func NewGrpcServer(server *grpc.Server, port ...int) *GrpcServer {
	return &GrpcServer{
		Base:   NewBase(GRPC, port...),
		server: server,
	}
}

// GrpcServer grpc服务
type GrpcServer struct {
	Base
	server *grpc.Server // grpc服务
}

// BindConfig 绑定服务配置
func (s *GrpcServer) BindConfig(config *Config) {
	s.bindConfig(config)
}

func (s *GrpcServer) Start(_ context.Context) error {
	if s.running {
		return nil
	} else if s.port == 0 {
		return errorx.New("server port is invalid")
	}
	addr := fmt.Sprintf(":%d", s.port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return errorx.Wrap(err, "listen error "+addr)
	}
	go func() {
		// 启动服务（非阻塞）
		if err = s.server.Serve(listen); err != nil {
			s.logger.WithError(err).Fatal("start server failed")
		}
	}()
	s.logger.Info("start server success")
	s.running = true
	return nil
}

func (s *GrpcServer) Shutdown(_ context.Context) {
	if s.running {
		s.server.Stop()
		s.running = false
	}
	s.logger.Info("shutdown server success")
}
