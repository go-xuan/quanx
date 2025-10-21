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
func NewGinServer(bindRouter func(engine *gin.Engine), middleware ...gin.HandlerFunc) *GinServer {
	// 创建gin引擎
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware...)

	// 绑定路由
	bindRouter(engine)

	return &GinServer{
		engine: engine,
	}
}

// BindRouter gin路由绑定函数
type BindRouter func(engine *gin.Engine)

// GinServer gin服务
type GinServer struct {
	engine *gin.Engine  // gin引擎
	srv    *http.Server // http服务
}

func (s *GinServer) Run(config *Config) error {
	if port := config.GetPort(); port != 0 {
		if config.Debug {
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
			log.WithField("port", port).Info("http server running")
			if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("http server listen error: %s\n", err)
			}
		}()
	}
	return nil
}

func (s *GinServer) Shutdown(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Error(err)
	}
	log.Info(`http server shutdown`)
}
