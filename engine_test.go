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
		AddTaskBefore(StepInitConfig, "before_init_config", func(context.Context) error {
			fmt.Println("before init config")
			return nil
		}),
		AddTaskAfter(StepInitConfig, "after_init_config", func(context.Context) error {
			fmt.Println("after init config")
			return nil
		}),
		AddServer(
			serverx.NewGinServer(func(engine *gin.Engine) {
				engine.Use(ginx.Trace)
				engine.GET("/ping", func(c *gin.Context) {
					c.String(http.StatusOK, "pong")
					return
				})
			})),
	).RUN(t.Context())
}
