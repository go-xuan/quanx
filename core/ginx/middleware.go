package ginx

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	clientIpKey  = "CLIENT_IP"
	localIp      = "::1"
	localIpValue = "127.0.0.1"
)

// RequestLogFmt gin请求日志格式化
func RequestLogFmt(ctx *gin.Context) {
	start := time.Now()
	// 处理请求
	ctx.Next()
	// 日志格式
	logrus.Infof("[%3d][%4dms][%s][%-4s %s]",
		ctx.Writer.Status(),
		time.Now().Sub(start).Milliseconds(),
		ClientIP(ctx),
		ctx.Request.Method,
		ctx.Request.URL.Path,
	)
}

// CorrectIP 纠正客户端IP
func CorrectIP(ctx *gin.Context) {
	if clientIP := ctx.ClientIP(); clientIP == localIp {
		ctx.Set(clientIpKey, localIpValue)
	} else {
		ctx.Set(clientIpKey, clientIP)
	}
	ctx.Next()
	return
}

// ClientIP 获取客户端IP
func ClientIP(ctx *gin.Context) string {
	if clientIP, ok := ctx.Get(clientIpKey); ok {
		return clientIP.(string)
	} else if clientIP = ctx.ClientIP(); clientIP == localIp {
		return localIpValue
	}
	return ""
}

// GetCorrectIP 当前请求IP
func GetCorrectIP(ctx *gin.Context) string {
	return ctx.GetString("ip")
}
