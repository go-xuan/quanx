package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-xuan/typex"
)

const (
	SuccessCode      = 10000 // 请求成功
	FailedCode       = 99999 // 请求失败
	AuthFailedCode   = 10401 // 鉴权失败
	ParamErrorCode   = 10501 // 请求参数错误
	RequiredCode     = 10502 // 请求参数必填
	UploadFailedCode = 10601 // 上传失败
	ImportFailedCode = 10602 // 导入失败
	ExportFailedCode = 10603 // 导出失败
)

var CodeEnum = typex.NewEnum[int, string]()

func init() {
	CodeEnum.Add(SuccessCode, "success").
		Add(FailedCode, "failed").
		Add(AuthFailedCode, "auth failed").
		Add(ParamErrorCode, "parse param failed").
		Add(RequiredCode, "param required").
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

// JsonAndAbort 自定义响应体
func JsonAndAbort(ctx *gin.Context, code int, data any) {
	ctx.JSON(code, data)
	ctx.Abort()
}

// CustomResponse 自定义错误
func CustomResponse(ctx *gin.Context, response *Response) {
	var code int
	switch response.Code {
	case SuccessCode:
		code = http.StatusOK
	case AuthFailedCode:
		code = http.StatusForbidden
	case ParamErrorCode:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}
	JsonAndAbort(ctx, code, response)
}

// Success 请求成功
func Success(ctx *gin.Context, data any) {
	CustomResponse(ctx, NewResponse(SuccessCode, data))
}

// Error 请求失败
func Error(ctx *gin.Context, err error) {
	CustomResponse(ctx, NewResponse(FailedCode, err.Error()))
}

// ParamError 请求参数错误
func ParamError(ctx *gin.Context, err error) {
	CustomResponse(ctx, NewResponse(ParamErrorCode, err.Error()))
}

// Forbidden 鉴权失败
func Forbidden(ctx *gin.Context, err error) {
	CustomResponse(ctx, NewResponse(AuthFailedCode, err.Error()))
}
