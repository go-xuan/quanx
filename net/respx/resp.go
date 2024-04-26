package respx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 正常响应
type Response struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
	Data any    `json:"data"` // 响应数据
}

// 响应
func BuildResponse(ctx *gin.Context, data any, err error) {
	if err != nil {
		BuildError(ctx, err)
	} else {
		BuildSuccess(ctx, data)
	}
}

// 成功响应
func BuildSuccess(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, &Response{Code: Success.Code, Msg: Success.Msg, Data: data})
}

// 错误响应
func BuildError(ctx *gin.Context, err error) {
	log.Error(Error.Msg, err)
	ctx.JSON(http.StatusInternalServerError, &Response{Code: Error.Code, Msg: Error.Msg, Data: err.Error()})
}

// 自定义错误响应
func Custom(ctx *gin.Context, code int, msg string, data any) {
	log.Error(msg, data)
	ctx.JSON(http.StatusInternalServerError, &Response{Code: code, Msg: msg, Data: data})
}

// 异常响应
func Exception(ctx *gin.Context, enum Enum, data any) {
	log.Error(enum.Msg, data)
	ctx.JSON(http.StatusInternalServerError, &Response{Code: enum.Code, Msg: enum.Msg, Data: data})
}

// 参数必填
func Required(ctx *gin.Context, msg string) {
	Exception(ctx, RequiredErr, msg)
}
