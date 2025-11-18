package quanx

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-xuan/quanx/ginx"
	"github.com/go-xuan/quanx/serverx"
)

func TestEngineRun(t *testing.T) {
	NewEngine(
		InitServerConfig(&serverx.Config{
			Name: "quanx-test", // 服务名称
			Port: map[string]int{
				"http": 8080,
			},
		}),
		// 添加http服务器
		AddServer(httpServer()),

		AddTaskBefore(FlagInitConfig, "before_init_config", func(context.Context) error {
			fmt.Println("do something before init config")
			return nil
		}),
		AddTaskAfter(FlagInitConfig, "after_init_config", func(context.Context) error {
			fmt.Println("do something after init config")
			return nil
		}),
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
