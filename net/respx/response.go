package respx

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type CtxBuilder struct {
	ctx *gin.Context
}

func Ctx(ctx *gin.Context) *CtxBuilder {
	return &CtxBuilder{ctx}
}

func (b *CtxBuilder) Response(data any, err error) {
	if err != nil {
		Error(b.ctx, ErrorCode, err)
	} else {
		Success(b.ctx, data)
	}
}

func (b *CtxBuilder) File(filePath string, err error) {
	if err != nil {
		Error(b.ctx, ErrorCode, err)
	} else {
		b.ctx.File(filePath)
	}
}

func (b *CtxBuilder) Render(contentType string, data []byte) {
	b.ctx.Render(http.StatusOK, render.Data{
		ContentType: contentType,
		Data:        data,
	})
}

func (b *CtxBuilder) Required(msg string) {
	b.ctx.JSON(http.StatusInternalServerError, RequiredCode.Response(msg))
}

// Response 正常响应
type Response struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
	Data any    `json:"data"` // 响应数据
}

// Success 成功响应
func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, SuccessCode.Response(data))
}

// Error 失败响应
func Error(ctx *gin.Context, enum *Enum, err error) {
	log.Error("请求失败：", err)
	ctx.JSON(http.StatusInternalServerError, enum.Response(err))
}

// Custom 自定义失败响应
func Custom(ctx *gin.Context, code int, msg string, data any) {
	ctx.JSON(http.StatusInternalServerError, &Response{Code: code, Msg: msg, Data: data})
}

// ParamError 请求参数错误
func ParamError(ctx *gin.Context, err error) {
	log.Error("参数错误：", err)
	Error(ctx, ParamErrorCode, err)
}

// Forbidden 鉴权失败
func Forbidden(ctx *gin.Context, err error) {
	log.Error("鉴权失败：", err)
	ctx.JSON(http.StatusForbidden, AuthFailedCode.Response(err))
}
