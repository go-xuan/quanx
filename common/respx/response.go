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

// 自定义响应
func BuildResponse(context *gin.Context, code int, message string, data interface{}) {
	response := Response{
		Code: code,
		Msg:  message,
		Data: data,
	}
	context.JSON(http.StatusInternalServerError, &response)
}

// 成功响应
func BuildSuccessResponse(context *gin.Context, data interface{}) {
	response := Response{
		Code: Success.Code,
		Msg:  Success.Msg,
		Data: data,
	}
	context.JSON(http.StatusOK, &response)
}

// 失败响应
func BuildFailResponse(context *gin.Context, data interface{}) {
	response := Response{
		Code: Fail.Code,
		Msg:  Fail.Msg,
		Data: data,
	}
	context.JSON(http.StatusBadRequest, &response)
}

// 错误响应
func BuildErrorResponse(context *gin.Context, data interface{}) {
	response := Response{
		Code: Error.Code,
		Msg:  Error.Msg,
		Data: data,
	}
	context.JSON(http.StatusInternalServerError, &response)
}

// 异常响应
func BuildExceptionResponse(context *gin.Context, res ErrEnum, data interface{}) {
	response := Response{
		Code: res.Code,
		Msg:  res.Msg,
		Data: data,
	}
	context.JSON(http.StatusInternalServerError, &response)
}

// 文件响应
func BuildFileResponse(context *gin.Context, filePath string) {
	context.File(filePath)
}
