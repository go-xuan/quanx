package ginx

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	sessionUserKey = "X_SESSION_USER"
	clientIpKey    = "X_CLIENT_IP"
	traceIdKey     = "X_TRACE_ID"
)

// DefaultEngine 创建gin引擎
func DefaultEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery()) // 恢复中间件
	return engine
}

// GetTraceId 获取traceId
func GetTraceId(ctx *gin.Context) string {
	if traceId, ok := ctx.Get(traceIdKey); ok {
		return traceId.(string)
	}
	return ""
}

// GetClientIP 获取客户端IP
func GetClientIP(ctx *gin.Context) string {
	var clientIp string
	if ip, ok := ctx.Get(clientIpKey); ok {
		clientIp = ip.(string)
	} else if clientIp = ctx.ClientIP(); clientIp == "::1" {
		clientIp = "127.0.0.1"
		ctx.Set(clientIpKey, clientIp)
	}
	return clientIp
}

// GetLogger 获取日志包装
func GetLogger(ctx *gin.Context) *log.Entry {
	entry := log.WithContext(ctx).
		WithField("trace_id", GetTraceId(ctx)).
		WithField("client_ip", GetClientIP(ctx))
	if user := GetSessionUser(ctx); user != nil {
		entry = entry.WithField("user_id", user.GetUserId())
	}
	return entry
}
