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
func NewGrpcServer(serviceDesc *grpc.ServiceDesc, serviceImpl interface{}, opt ...grpc.ServerOption) *GrpcServer {
	// 创建grpc服务
	server := grpc.NewServer(opt...)

	// 注册服务
	server.RegisterService(serviceDesc, serviceImpl)

	return &GrpcServer{
		server: server,
	}
}

// GrpcServer grpc服务
type GrpcServer struct {
	server *grpc.Server
}

func (s *GrpcServer) Run(config *Config) error {
	if port := config.Grpc; port != 0 {
		address := fmt.Sprintf(":%d", port)
		listen, err := net.Listen("tcp", address)
		if err != nil {
			return errorx.Wrap(err, "listen error")
		}
		go s.server.Serve(listen)
		log.WithField("port", port).Info("grpc server running")
	}
	return nil
}

func (s *GrpcServer) Shutdown(_ context.Context) {
	s.server.Stop()
	log.Info(`grpc server shutdown`)
}
