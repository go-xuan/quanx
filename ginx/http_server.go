package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HttpServer 创建http服务器
func HttpServer(options ...EngineOption) *http.Server {
	engine := DefaultEngine()
	for _, option := range options {
		option(engine)
	}
	return &http.Server{
		Handler: engine,
	}
}

// EngineOption engine选项
type EngineOption func(engine *gin.Engine)

// SetDebugMode 设置gin为调试模式
func SetDebugMode(_ *gin.Engine) {
	gin.SetMode(gin.DebugMode)
}

// AddMiddleware 添加中间件
func AddMiddleware(middlewares ...gin.HandlerFunc) EngineOption {
	return func(engine *gin.Engine) {
		engine.Use(middlewares...)
	}
}
