package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewHttpServer 创建http服务器
func NewHttpServer(options ...EngineOption) *http.Server {
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
