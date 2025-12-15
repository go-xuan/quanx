package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/go-xuan/typex"
)

const (
	SuccessCode      = 10000
	FailedCode       = 99999
	AuthFailedCode   = 10401
	ParamErrorCode   = 10501
	RequiredCode     = 10502
	UploadFailedCode = 10601
	ImportFailedCode = 10602
	ExportFailedCode = 10603
)

var CodeEnum = typex.NewEnum[int, string]()

func init() {
	CodeEnum.Add(SuccessCode, "success").
		Add(FailedCode, "failed").
		Add(AuthFailedCode, "auth failed").
		Add(ParamErrorCode, "request parameter error").
		Add(RequiredCode, "request parameter required").
		Add(UploadFailedCode, "upload failed").
		Add(ImportFailedCode, "import failed").
		Add(ExportFailedCode, "export failed")
}

// Response 响应数据
type Response struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
	Data any    `json:"data"` // 响应数据
}

// NewResponse 初始化响应数据
func NewResponse(code int, data any) *Response {
	return &Response{
		Code: code,
		Msg:  CodeEnum.Get(code),
		Data: data,
	}
}

// Success 请求成功
func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, NewResponse(SuccessCode, data))
}

// Error 请求失败
func Error(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, NewResponse(FailedCode, err.Error()))
}

// ParamError 请求参数错误
func ParamError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, NewResponse(ParamErrorCode, err.Error()))
}

// Forbidden 鉴权失败
func Forbidden(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusForbidden, NewResponse(AuthFailedCode, err.Error()))
}

// Custom 自定义响应体
func Custom(ctx *gin.Context, httpCode int, data *Response) {
	ctx.JSON(httpCode, data)
}

// CustomError 自定义错误
func CustomError(ctx *gin.Context, data *Response) {
	code := http.StatusInternalServerError
	switch data.Code {
	case AuthFailedCode:
		code = http.StatusForbidden
	case ParamErrorCode:
		code = http.StatusBadRequest
	}
	ctx.JSON(code, data)
}

// File 文件输出
func File(ctx *gin.Context, filepath string) {
	ctx.File(filepath)
}

func Render(ctx *gin.Context, data render.Data) {
	ctx.Render(http.StatusOK, data)
}
