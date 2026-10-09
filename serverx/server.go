package serverx

import (
	"context"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

const (
	HTTP  = "http"  // http服务类型
	GRPC  = "grpc"  // grpc服务类型
	PPROF = "pprof" // pprof服务类型
)

// Server 服务
type Server interface {
	BindConfig(config *Config)       // 绑定服务配置
	Start(ctx context.Context) error // 启动服务
	Shutdown(ctx context.Context)    // 关闭服务
}

// NewBase 创建基础服务
func NewBase(category string, port ...int) Base {
	base := Base{category: category}
	if len(port) > 0 && port[0] > 0 {
		base.port = port[0]
	}
	return base
}

// Base 基础服务配置
type Base struct {
	name     string     // 服务名称
	category string     // 服务分类
	port     int        // 服务端口
	running  bool       // 服务运行标识
	logger   *log.Entry // 日志记录器
}

// 绑定服务配置
func (s *Base) bindConfig(config *Config) {
	if config == nil {
		return
	}
	// 绑定服务名称
	if name := config.Name; s.name == "" && name != "" {
		s.name = name
	}
	// 绑定端口
	if port := config.Port[s.category]; s.port == 0 && port > 0 {
		s.port = port
	}
	s.logger = log.WithFields(log.Fields{
		"name":     s.name,
		"category": s.category,
		"port":     s.port,
	})
}

// Start 启动服务
func Start(ctx context.Context, config *Config, servers ...Server) error {
	for _, server := range servers {
		server.BindConfig(config)
		if err := server.Start(ctx); err != nil {
			return errorx.Wrap(err, "server start error")
		}
	}
	// 注册服务
	if err := config.RegisterServer(); err != nil {
		return errorx.Wrap(err, "server register failed")
	}
	return nil
}

// Shutdown 关闭服务
func Shutdown(ctx context.Context, servers ...Server) {
	for _, server := range servers {
		server.Shutdown(ctx)
	}
}
