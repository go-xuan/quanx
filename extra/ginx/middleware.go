package ginx

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// DefaultLogFormatter gin请求日志格式化
func DefaultLogFormatter(ctx *gin.Context) {
	start := time.Now()
	// 处理请求
	ctx.Next()
	// 日志格式化
	Log(ctx).Infof("[%3d][%dms][%-4s %s]",
		ctx.Writer.Status(),
		time.Since(start).Milliseconds(),
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
	Log(ctx).WithField("method", ctx.Request.Method).
		WithField("url", ctx.Request.URL.Path).
		WithField("status", ctx.Writer.Status()).
		WithField("duration", time.Since(start).Milliseconds()).
		Info()
}

// Log 日志包装
func Log(ctx *gin.Context) *log.Entry {
	entry := log.WithField("traceId", TraceId(ctx)).WithField("clientIp", ClientIP(ctx))
	if user := GetSessionUser(ctx); user != nil {
		entry = entry.WithField("userId", user.UserId())
	}
	return entry
}

// ClientIP 获取客户端IP
func ClientIP(ctx *gin.Context) string {
	var clientIp string
	if ip, ok := ctx.Get(clientIpKey); ok {
		clientIp = ip.(string)
	} else if clientIp = ctx.ClientIP(); clientIp == "::1" {
		clientIp = "127.0.0.1"
		ctx.Set(clientIpKey, clientIp)
	}
	return clientIp
}

// Trace traceId
func Trace(ctx *gin.Context) {
	ctx.Set(traceIdKey, uuid.NewString())
	ctx.Next()
}

// TraceId 获取traceId
func TraceId(ctx *gin.Context) string {
	if traceId, ok := ctx.Get(traceIdKey); ok {
		return traceId.(string)
	}
	return ""
}
