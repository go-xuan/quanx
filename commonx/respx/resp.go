package respx

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 正常响应
type Response struct {
	Code int         `json:"code"` // 响应状态码
	Msg  string      `json:"msg"`  // 响应消息
	Data interface{} `json:"data"` // 响应数据
}

// 响应
func BuildResponse(ctx *gin.Context, data interface{}, err error) {
	if err != nil {
		BuildError(ctx, err)
	} else {
		BuildSuccess(ctx, data)
	}
}

// 成功响应
func BuildSuccess(ctx *gin.Context, data interface{}) {
	response := &Response{
		Code: Success.Code,
		Msg:  Success.Msg,
		Data: data,
	}
	ctx.JSON(http.StatusOK, response)
}

// 错误响应
func BuildError(ctx *gin.Context, err error) {
	response := &Response{
		Code: Error.Code,
		Msg:  Error.Msg,
		Data: err.Error(),
	}
	log.Error("server error ： ", err)
	ctx.JSON(http.StatusInternalServerError, response)
}

// 自定义错误响应
func Custom(ctx *gin.Context, code int, msg string, data interface{}) {
	response := &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	log.Error(msg, data)
	ctx.JSON(http.StatusInternalServerError, response)
}

// 异常响应
func Exception(ctx *gin.Context, enum Enum, data interface{}) {
	response := &Response{
		Code: enum.Code,
		Msg:  enum.Msg,
		Data: data,
	}
	log.Error(enum.Msg, data)
	ctx.JSON(http.StatusInternalServerError, response)
}

// 参数必填
func Required(ctx *gin.Context, msg string) {
	Exception(ctx, RequiredErr, msg)
}
