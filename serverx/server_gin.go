package serverx

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// NewGinServer 创建gin服务
func NewGinServer(prepare func(*gin.Engine)) *GinServer {
	// 创建gin引擎
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	return &GinServer{
		engine:  engine,
		prepare: prepare,
	}
}

// GinServer gin服务
type GinServer struct {
	config  *Config           // 服务配置
	engine  *gin.Engine       // gin引擎
	prepare func(*gin.Engine) // gin引擎准备
	srv     *http.Server      // http服务
	running bool              // 服务运行标识
}

func (s *GinServer) Run(config *Config) error {
	if port := config.GetPort(); port != 0 {
		// 配置gin引擎
		s.prepare(s.engine)

		if config.Debug {
			gin.SetMode(gin.DebugMode)
			for _, info := range s.engine.Routes() {
				log.WithField("method", info.Method).
					WithField("path", info.Path).
					WithField("handler", info.Handler).
					Debug("http route info")
			}
		}
		s.srv = &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: s.engine,
		}
		go func() {
			// 启动服务（非阻塞）
			log.WithField("name", config.GetName()).
				WithField("port", config.GetPort()).
				WithField("debug", config.Debug).
				Info("http server running")
			if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("http server listen error: %s\n", err)
			}
		}()
		s.config = config
		s.running = true
	}
	return nil
}

func (s *GinServer) Shutdown(ctx context.Context) {
	if s.running {
		if err := s.srv.Shutdown(ctx); err != nil {
			log.Error(err)
		}
		s.running = false
	}
	log.Info(`http server shutdown`)
}

func (s *GinServer) IsRunning() bool {
	return s.running
}
