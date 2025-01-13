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

// DefaultLogFormatter gin请求日志格式化
func DefaultLogFormatter(ctx *gin.Context) {
	start := time.Now()
	// 处理请求
	ctx.Next()
	// 日志格式化
	logrus.Infof("[%3d][%4dms][%s][%-4s %s]",
		ctx.Writer.Status(),
		time.Now().Sub(start).Milliseconds(),
		ClientIP(ctx),
		ctx.Request.Method,
		ctx.Request.URL.Path,
	)
}

// JsonLogFormatter gin请求日志格式化
func JsonLogFormatter(ctx *gin.Context) {
	start := time.Now()
	// 处理请求
	ctx.Next()
	// 日志格式化
	logrus.WithField("method", ctx.Request.Method).
		WithField("url", ctx.Request.URL.Path).
		WithField("status", ctx.Writer.Status()).
		WithField("duration", time.Since(start).String()).
		WithField("clientIp", ClientIP(ctx)).
		Info()
}

// ClientIP 获取客户端IP
func ClientIP(ctx *gin.Context) string {
	if ip, ok := ctx.Get(clientIpKey); ok {
		return ip.(string)
	} else if clientIp := ctx.ClientIP(); clientIp == localIp {
		ctx.Set(clientIpKey, localIpValue)
		return localIpValue
	} else {
		return clientIp
	}
}
