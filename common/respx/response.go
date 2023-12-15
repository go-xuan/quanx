package respx

import (
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
	ctx.JSON(http.StatusInternalServerError, response)
}

// 自定义错误响应
func BuildCustomErr(ctx *gin.Context, code int, msg string, data interface{}) {
	response := &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	ctx.JSON(http.StatusInternalServerError, response)
}

// 异常响应
func BuildException(ctx *gin.Context, res ErrorEnum, data interface{}) {
	response := &Response{
		Code: res.Code,
		Msg:  res.Msg,
		Data: data,
	}
	ctx.JSON(http.StatusInternalServerError, response)
}

// 文件响应
func BuildFile(ctx *gin.Context, filePath string) {
	ctx.File(filePath)
}
