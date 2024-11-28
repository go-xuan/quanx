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

func NewResponseData(code int, data any) *ResponseData {
	return &ResponseData{
		Code: code,
		Msg:  CodeMsgEnum.Get(code),
		Data: data,
	}
}

func Response(ctx *gin.Context, data any, err error) {
	if err != nil {
		Error(ctx, err.Error())
	} else {
		Success(ctx, data)
	}
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, NewResponseData(SuccessCode, data))
}

func Error(ctx *gin.Context, data any) {
	log.Errorf("[%s]请求失败：%v", ctx.Request.URL.Path, data)
	ctx.JSON(http.StatusInternalServerError, NewResponseData(FailedCode, data))
}

func ParamError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, NewResponseData(ParamErrorCode, err.Error()))
}

func Forbidden(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusForbidden, NewResponseData(AuthFailedCode, err.Error()))
}

func Custom(ctx *gin.Context, httpCode int, data *ResponseData) {
	ctx.JSON(httpCode, data)
}

func CustomError(ctx *gin.Context, data *ResponseData) {
	ctx.JSON(http.StatusInternalServerError, data)
}

func File(ctx *gin.Context, filepath string) {
	ctx.File(filepath)
}

func Render(ctx *gin.Context, data render.Data) {
	ctx.Render(http.StatusOK, data)
}
