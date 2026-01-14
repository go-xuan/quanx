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
		Base:   NewBase(HTTP, port...),
		server: server,
	}
}

// NewPprofServer 创建pprof服务
func NewPprofServer(port ...int) *HttpServer {
	// 需要将 _ "net/http/pprof" 添加到项目的main.go文件的import语句中
	// 使用默认的 http.DefaultServeMux，会自动注册pprof相关路由
	server := &http.Server{
		Handler: http.DefaultServeMux,
	}
	return &HttpServer{
		Base:   NewBase(PPROF, port...),
		server: server,
	}
}

// HttpServer http服务
type HttpServer struct {
	Base                // 基础服务配置
	server *http.Server // http服务
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
			s.logger.WithError(err).Fatal("start server failed")
		}
	}()
	s.logger.Info("start server success")
	s.running = true
	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) {
	if s.running {
		if err := s.server.Shutdown(ctx); err != nil {
			s.logger.WithError(err).Error("shutdown server failed")
			return
		}
		s.running = false
	}
	s.logger.Info("shutdown server success")
}
