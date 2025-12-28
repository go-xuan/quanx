package ginx

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Cors 跨域处理
func Cors(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-type, Authorization, Origin, Accept, User-Agent")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}

// Trace 为gin上下文添加traceId
func Trace(ctx *gin.Context) {
	ctx.Set(traceIdKey, uuid.NewString())
	ctx.Next()
}

// LogFormatter gin请求日志格式化
func LogFormatter(ctx *gin.Context) {
	start := time.Now()
	// 处理请求
	ctx.Next()
	// 日志格式化
	GetLogger(ctx).WithField("method", ctx.Request.Method).
		WithField("url", ctx.Request.URL.String()).
		WithField("status", ctx.Writer.Status()).
		WithField("duration", time.Since(start).String()).
		Info("gin request finished")
}

// AdvanceBindJSON 提前绑定JSON数据，一般用于在前置中间件中绑定JSON数据，后续中间件可以使用绑定的数据
func AdvanceBindJSON(ctx *gin.Context, data interface{}) error {
	// 提前拿出请求体
	body, _ := io.ReadAll(ctx.Request.Body)
	defer SetCtxBody(ctx, body)

	SetCtxBody(ctx, body)
	// 绑定数据
	if err := ctx.ShouldBindJSON(data); err != nil {
		return err
	}
	return nil
}

// SetCtxBody 设置gin请求体
func SetCtxBody(ctx *gin.Context, body []byte) {
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
}
