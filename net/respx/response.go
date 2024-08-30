package respx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	log "github.com/sirupsen/logrus"
)

// Response 正常响应
type Response struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
	Data any    `json:"data"` // 响应数据
}

type ResponseBuilder struct {
	ctx *gin.Context
	f   func()
}

func Ctx(ctx *gin.Context) *ResponseBuilder {
	return &ResponseBuilder{ctx: ctx}
}

func WhenSuccess(ctx *gin.Context, f func()) *ResponseBuilder {
	return &ResponseBuilder{ctx: ctx, f: f}
}

func (b *ResponseBuilder) Success(data any) {
	if b.f != nil {
		b.f()
	}
	b.ctx.JSON(http.StatusOK, SuccessEnum.Response(data))
}

func (b *ResponseBuilder) Error(err error) {
	log.Error("请求失败：", err)
	b.ctx.JSON(http.StatusInternalServerError, ErrorEnum.Response(err.Error()))
}

func (b *ResponseBuilder) EnumError(enum *Enum, err error) {
	log.Error("请求失败：", err)
	b.ctx.JSON(http.StatusInternalServerError, enum.Response(err.Error()))
}

func (b *ResponseBuilder) CustomError(code int, msg string, data any) {
	b.ctx.JSON(http.StatusInternalServerError, &Response{Code: code, Msg: msg, Data: data})
}

func (b *ResponseBuilder) ParamError(err error) {
	log.Error("参数错误：", err)
	b.ctx.JSON(http.StatusInternalServerError, ParamErrorEnum.Response(err.Error()))
}

func (b *ResponseBuilder) Forbidden(err error) {
	log.Error("鉴权失败：", err)
	b.ctx.JSON(http.StatusForbidden, AuthFailedEnum.Response(err.Error()))
}

func (b *ResponseBuilder) Response(data any, err error) {
	if err != nil {
		b.Error(err)
	} else {
		b.Success(data)
	}
}

func (b *ResponseBuilder) File(filePath string) {
	b.ctx.File(filePath)
}

func (b *ResponseBuilder) Render(contentType string, data []byte) {
	b.ctx.Render(http.StatusOK, render.Data{
		ContentType: contentType,
		Data:        data,
	})
}
