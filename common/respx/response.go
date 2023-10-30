package respx

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

// 自定义响应
func BuildCustom(ctx *gin.Context, code int, message string, data interface{}) {
	response := Response{
		Code: code,
		Msg:  message,
		Data: data,
	}
	ctx.JSON(http.StatusInternalServerError, &response)
}

// 成功响应
func BuildSuccess(ctx *gin.Context, data interface{}) {
	response := Response{
		Code: Success.Code,
		Msg:  Success.Msg,
		Data: data,
	}
	ctx.JSON(http.StatusOK, &response)
}

// 错误响应
func BuildError(ctx *gin.Context, err error) {
	response := Response{
		Code: Error.Code,
		Msg:  Error.Msg,
		Data: err.Error(),
	}
	ctx.JSON(http.StatusInternalServerError, &response)
}

// 异常响应
func BuildException(ctx *gin.Context, res ErrEnum, data interface{}) {
	response := Response{
		Code: res.Code,
		Msg:  res.Msg,
		Data: data,
	}
	ctx.JSON(http.StatusInternalServerError, &response)
}

// 文件响应
func BuildFileResponse(ctx *gin.Context, filePath string) {
	ctx.File(filePath)
}
