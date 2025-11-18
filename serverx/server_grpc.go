package serverx

import (
	"context"
	"fmt"
	"net"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// NewGrpcServer 创建grpc服务
func NewGrpcServer(server *grpc.Server, port ...int) *GrpcServer {
	grpcServer := GrpcServer{
		server: server,
	}
	if len(port) > 0 && port[0] > 0 {
		grpcServer.port = port[0]
	}
	return &grpcServer
}

// GrpcServer grpc服务
type GrpcServer struct {
	running bool // 服务运行标识
	name    string
	port    int          // grpc服务端口
	server  *grpc.Server // grpc服务
}

// IsRunning 是否运行中
func (s *GrpcServer) IsRunning() bool {
	return s.running
}

// BindConfig 绑定服务配置
func (s *GrpcServer) BindConfig(config *Config) {
	if config == nil {
		return
	}
	// 绑定服务名称
	if name := config.Name; s.name == "" && name != "" {
		s.name = name
	}

	// 绑定端口
	if port := config.Port["grpc"]; s.port == 0 && port > 0 {
		s.port = port
	}
}

func (s *GrpcServer) Start(ctx context.Context) error {
	// 检查服务是否已运行或配置是否存在
	if s.running {
		return errorx.New("grpc server is running")
	} else if s.port == 0 {
		return errorx.New("grpc server port is not set")
	}
	logger := log.WithField("type", "grpc").WithField("name", s.name).WithField("port", s.port)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return errorx.Wrap(err, "listen error")
	}
	go func() {
		// 启动服务（非阻塞）
		if err = s.server.Serve(listen); err != nil {
			logger.WithError(err).Fatal("server run failed")
		}
	}()
	logger.Info("server run success")
	s.running = true
	return nil
}

func (s *GrpcServer) Shutdown(_ context.Context) {
	logger := log.WithField("type", "grpc").WithField("name", s.name).WithField("port", s.port)
	if !s.running {
		logger.Warn("server not running")
		return
	}
	s.server.Stop()
	s.running = false
	logger.Info("server shutdown success")
}
