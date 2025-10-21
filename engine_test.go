package quanx

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-xuan/quanx/ginx"
	"github.com/go-xuan/quanx/serverx"
)

func TestEngineRun(t *testing.T) {
	// 初始化Engine
	e := NewEngine(
		AddServer(serverx.NewGinServer(func(engine *gin.Engine) {
			engine.GET("/ping", func(c *gin.Context) {
				c.String(http.StatusOK, "pong")
				return
			})
		}, ginx.Trace)))
	// 启动服务
	e.RUN(t.Context())
}
