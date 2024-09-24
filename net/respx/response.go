package respx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	log "github.com/sirupsen/logrus"
)

var (
	SuccessResp      = Response{Code: 10000, Msg: "Success"}
	FailedResp       = Response{Code: 10001, Msg: "failed"}
	AuthFailedResp   = Response{Code: 10401, Msg: "auth failed"}
	ParamErrorResp   = Response{Code: 10501, Msg: "request parameter error"}
	RequiredResp     = Response{Code: 10502, Msg: "request parameter required"}
	UploadFailedResp = Response{Code: 10601, Msg: "upload failed"}
	ImportFailedResp = Response{Code: 10601, Msg: "import failed"}
	ExportFailedResp = Response{Code: 10602, Msg: "export failed"}
)

// Response 响应
type Response struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
	Data any    `json:"data"` // 响应数据
}

// SetData 设置响应体
func (e Response) SetData(data any) *Response {
	return &Response{
		Code: e.Code,
		Msg:  e.Msg,
		Data: data,
	}
}

type Builder struct {
	ctx *gin.Context
}

func Ctx(ctx *gin.Context) *Builder {
	return &Builder{
		ctx: ctx,
	}
}

func (b *Builder) Response(data any, err error) {
	if err != nil {
		log.Error("请求失败：", err)
		b.Failed(err)
	} else {
		b.Success(data)
	}
	b.ctx.Next()
}

func (b *Builder) Success(data any) {
	b.ctx.JSON(http.StatusOK, SuccessResp.SetData(data))
}

func (b *Builder) Failed(err error) {
	b.ctx.JSON(http.StatusInternalServerError, FailedResp.SetData(err.Error()))
}

func (b *Builder) RespError(resp *Response) {
	b.ctx.JSON(http.StatusInternalServerError, resp)
}

func (b *Builder) ParamError(err error) {
	b.ctx.JSON(http.StatusInternalServerError, ParamErrorResp.SetData(err.Error()))
}

func (b *Builder) Forbidden(err error) {
	b.ctx.JSON(http.StatusForbidden, AuthFailedResp.SetData(err.Error()))
}

func (b *Builder) Custom(httpCode, code int, msg string, data any) {
	b.ctx.JSON(httpCode,
		&Response{
			Code: code,
			Msg:  msg,
			Data: data,
		})
}

func (b *Builder) CustomError(code int, msg string, data any) {
	b.Custom(http.StatusInternalServerError, code, msg, data)
}

func (b *Builder) File(filePath string) {
	b.ctx.File(filePath)
}

func (b *Builder) Render(contentType string, data []byte) {
	b.ctx.Render(http.StatusOK, render.Data{
		ContentType: contentType,
		Data:        data,
	})
}
