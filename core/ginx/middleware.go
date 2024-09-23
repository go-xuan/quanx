package ginx

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RequestLogFmt gin请求日志格式化
func RequestLogFmt(ctx *gin.Context) {
	start := time.Now()
	// 处理请求
	ctx.Next()
	var ip string
	if ipv, ok := ctx.Get("ip"); ok {
		ip = ipv.(string)
	} else if ip = ctx.ClientIP(); ip == "::1" {
		ip = "127.0.0.1"
	}
	// 日志格式
	logrus.Infof("[%3d][%4dms][%s][%-4s %s]",
		ctx.Writer.Status(),
		time.Now().Sub(start).Milliseconds(),
		ip,
		ctx.Request.Method,
		ctx.Request.URL.Path,
	)
}

// CheckIP 校验请求IP
func CheckIP(ctx *gin.Context) {
	var ip string
	if ip = ctx.ClientIP(); ip == "::1" {
		ip = "127.0.0.1"
	}
	ctx.Set("ip", ip)
	ctx.Next()
	return
}

// GetCorrectIP 当前请求IP
func GetCorrectIP(ctx *gin.Context) string {
	return ctx.GetString("ip")
}
