package appx

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-xuan/quanx/ginx"
	"github.com/go-xuan/quanx/serverx"
)

func TestEngine(t *testing.T) {
	NewEngine(
		// 设置服务配置
		SetServerConfig(&serverx.Config{
			Name: "quanx-test", // 服务名称
			Host: "localhost",  // host, 为空时默认获取本地IP
			Port: map[string]int{
				serverx.HTTP: 8080,
			},
		}),
		// 添加服务
		AddServer(
			// 添加http服务
			serverx.NewHttpServer(ginx.HttpServer(
				// 开启调试模式
				ginx.SetDebugMode,
				// 绑定路由
				func(engine *gin.Engine) {
					engine.GET("/ping", func(ctx *gin.Context) {
						ctx.String(http.StatusOK, "pong")
						return
					})
				},
			), 8081),
		),
	).RUN(t.Context())
}
