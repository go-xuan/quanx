package serverx

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-xuan/utilx/errorx"
)

// NewHttpServer 创建http服务
func NewHttpServer(server *http.Server, port ...int) *HttpServer {
	return &HttpServer{
		BaseServer: NewBaseServer(HTTP, port...),
		server:     server,
	}
}

// NewPprofServer 创建pprof服务
func NewPprofServer(port ...int) *HttpServer {
	return &HttpServer{
		BaseServer: NewBaseServer(PPROF, port...),
		server: &http.Server{
			// 使用默认的 http.DefaultServeMux，当import pprof包时，会自动注册pprof路由
			Handler: http.DefaultServeMux,
		},
	}
}

// HttpServer http服务
type HttpServer struct {
	BaseServer              // 基础服务配置
	server     *http.Server // http服务
}

// IsRunning 是否运行中
func (s *HttpServer) IsRunning() bool {
	return s.running
}

// BindConfig 绑定服务配置
func (s *HttpServer) BindConfig(config *Config) {
	s.bindConfig(config)
}

func (s *HttpServer) Start(_ context.Context) error {
	if s.running {
		return nil
	} else if s.port == 0 {
		return errorx.New("server port is invalid")
	}
	go func() {
		// 启动服务（非阻塞）
		s.server.Addr = fmt.Sprintf(":%d", s.port)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.WithError(err).Fatal("server run failed")
		}
	}()
	s.logger.Info("server start success")
	s.running = true
	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) {
	if s.running {
		if err := s.server.Shutdown(ctx); err != nil {
			s.logger.WithError(err).Error("server shutdown failed")
			return
		}
		s.running = false
	}
	s.logger.Info("server shutdown success")
}
