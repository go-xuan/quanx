package respx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/types/enumx"
)

const (
	SuccessCode      = 10000
	FailedCode       = 10001
	AuthFailedCode   = 10401
	ParamErrorCode   = 10501
	RequiredCode     = 10502
	UploadFailedCode = 10601
	ImportFailedCode = 10601
	ExportFailedCode = 10602
)

var CodeMsgEnum = enumx.NewIntEnum[string]()

func init() {
	CodeMsgEnum.Add(SuccessCode, "success").
		Add(FailedCode, "failed").
		Add(AuthFailedCode, "auth failed").
		Add(ParamErrorCode, "request parameter error").
		Add(RequiredCode, "request parameter required").
		Add(UploadFailedCode, "upload failed").
		Add(ImportFailedCode, "import failed").
		Add(ExportFailedCode, "export failed")
}

// ResponseData 响应数据
type ResponseData struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
	Data any    `json:"data"` // 响应数据
}

// NewResponseData 初始化响应数据
func NewResponseData(code int, data any) *ResponseData {
	return &ResponseData{
		Code: code,
		Msg:  CodeMsgEnum.Get(code),
		Data: data,
	}
}

// Success 请求成功
func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, NewResponseData(SuccessCode, data))
	log.WithContext(ctx).WithField("url", ctx.Request.URL.Path).
		WithField("code", http.StatusOK).WithField("data", data).
		Info("request success")
}

// Error 请求失败
func Error(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusInternalServerError, NewResponseData(FailedCode, data))
	log.WithContext(ctx).WithField("url", ctx.Request.URL.Path).
		WithField("code", http.StatusInternalServerError).WithField("data", data).
		Error("request error")
}

// ParamError 请求参数错误
func ParamError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, NewResponseData(ParamErrorCode, err.Error()))
	log.WithContext(ctx).WithField("url", ctx.Request.URL.Path).
		WithField("code", http.StatusBadRequest).WithError(err).
		Error("request parameter error")
}

// Forbidden 鉴权失败
func Forbidden(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusForbidden, NewResponseData(AuthFailedCode, err.Error()))
	log.WithContext(ctx).WithField("url", ctx.Request.URL.Path).
		WithField("code", http.StatusForbidden).WithError(err).
		Error("request auth validate error")
}

// Custom 自定义响应体
func Custom(ctx *gin.Context, code int, data *ResponseData) {
	ctx.JSON(code, data)
	log.WithContext(ctx).WithField("url", ctx.Request.URL.Path).
		WithField("code", code).WithField("data", data).
		Info("request custom")
}

// CustomError 自定义错误
func CustomError(ctx *gin.Context, data *ResponseData) {
	ctx.JSON(http.StatusInternalServerError, data)
	log.WithContext(ctx).WithField("url", ctx.Request.URL.Path).
		WithField("code", http.StatusInternalServerError).WithField("data", data).
		Error("request custom error")
}

// File 文件输出
func File(ctx *gin.Context, filepath string) {
	ctx.File(filepath)
	log.WithContext(ctx).WithField("url", ctx.Request.URL.Path).
		WithField("filepath", filepath).
		Info("request output file")
}

func Render(ctx *gin.Context, data render.Data) {
	ctx.Render(http.StatusOK, data)
	log.WithContext(ctx).WithField("url", ctx.Request.URL.Path).
		WithField("code", http.StatusOK).WithField("data", data).
		Info("request render success")
}
