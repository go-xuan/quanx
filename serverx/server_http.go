package serverx

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

// NewHttpServer 创建http服务
func NewHttpServer(server *http.Server, port ...int) *HttpServer {
	httpServer := HttpServer{
		server: server,
	}
	if len(port) > 0 && port[0] > 0 {
		httpServer.port = port[0]
	}
	return &httpServer
}

// HttpServer http服务
type HttpServer struct {
	running bool         // 服务运行标识
	name    string       // 服务名称
	port    int          // http服务端口
	server  *http.Server // http服务
}

// IsRunning 是否运行中
func (s *HttpServer) IsRunning() bool {
	return s.running
}

// BindConfig 绑定服务配置
func (s *HttpServer) BindConfig(config *Config) {
	if config == nil {
		return
	}

	// 绑定服务名称
	if name := config.Name; s.name == "" && name != "" {
		s.name = name
	}

	// 绑定端口
	if port := config.Port["http"]; s.port == 0 && port > 0 {
		s.port = port
	}
}

func (s *HttpServer) Start(ctx context.Context) error {
	// 检查服务是否已运行或配置是否存在
	if s.running {
		return errorx.New("http server is running")
	} else if s.port == 0 {
		return errorx.New("http server port is not set")
	}
	logger := log.WithField("type", "http").WithField("name", s.name).WithField("port", s.port)
	go func() {
		// 启动服务（非阻塞）
		s.server.Addr = fmt.Sprintf(":%d", s.port)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Fatal("server run failed")
		}
	}()
	logger.Info("server run success")
	s.running = true
	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) {
	logger := log.WithField("type", "http").WithField("name", s.name).WithField("port", s.port)
	if !s.running {
		logger.Warn("server not running")
		return
	}
	s.running = false
	if err := s.server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("server shutdown failed")
		return
	}
	logger.Info("server shutdown success")

}
