package quanx

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestEngineRun(t *testing.T) {
	// 初始化Engine
	e := GetEngine()

	e.DoOption(AddGinRouter(func(group *gin.RouterGroup) {
		group.GET("/hello", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "hello world")
		})
	}))

	// 启动服务
	e.RUN(t.Context())
}
