package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewHttpServer 创建http服务器
func NewHttpServer(builders ...EngineBuilder) *http.Server {
	engine := DefaultEngine()
	// 构造引擎
	for _, builder := range builders {
		builder(engine)
	}
	return &http.Server{
		Handler: engine,
	}
}

// EngineBuilder engine构造方法
type EngineBuilder func(engine *gin.Engine)

// SetDebugMode 设置gin为调试模式
func SetDebugMode(_ *gin.Engine) {
	gin.SetMode(gin.DebugMode)
}
