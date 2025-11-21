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
		SetServerConfig(serverConfig()), // 设置服务配置
		AddServer(httpServer()),         // 添加http服务
	).RUN(t.Context())
}

// 创建http服务
func httpServer() *serverx.HttpServer {
	return serverx.NewHttpServer(ginx.NewHttpServer(
		ginx.SetDebugMode, // 开启调试模式
		bindRouter,        // 绑定路由
	))
}

// 绑定路由
func bindRouter(engine *gin.Engine) {
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
		return
	})
}

func serverConfig() *serverx.Config {
	return &serverx.Config{
		Name: "quanx-test", // 服务名称
		Port: map[string]int{
			serverx.HTTP: 8080,
		},
	}
}
