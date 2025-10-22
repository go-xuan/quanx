package ginx

import "github.com/gin-gonic/gin"

const (
	userCookieKey  = "GIN_COOKIE_USER"
	sessionUserKey = "GIN_SESSION_USER"
	clientIpKey    = "GIN_CLIENT_IP"
	traceIdKey     = "GIN_TRACE_ID"
)

// PrepareFunc gin引擎准备函数
type PrepareFunc func(*gin.Engine)

// TraceId 获取traceId
func TraceId(ctx *gin.Context) string {
	if traceId, ok := ctx.Get(traceIdKey); ok {
		return traceId.(string)
	}
	return ""
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
