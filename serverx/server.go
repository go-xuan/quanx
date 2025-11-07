package serverx

import (
	"context"
)

// Server 服务
type Server interface {
	BindConfig(config *Config)       // 绑定服务配置
	Start(ctx context.Context) error // 启动服务
	Shutdown(ctx context.Context)    // 关闭服务
	IsRunning() bool                 // 是否运行中
}

// StartAll 启动服务
func StartAll(ctx context.Context, config *Config, servers ...Server) error {
	for _, server := range servers {
		server.BindConfig(config)
		if err := server.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

// ShutdownAll 关闭服务
func ShutdownAll(ctx context.Context, servers ...Server) {
	for _, server := range servers {
		server.Shutdown(ctx)
	}
}
