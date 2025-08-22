package ginx

import (
	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	Path string
}

func (h *BaseHandler) RelativePath() string {
	return h.Path
}

func (h *BaseHandler) Success(ctx *gin.Context, data any) {
	Success(ctx, data)
}

func (h *BaseHandler) Error(ctx *gin.Context, data error) {
	Error(ctx, data)
}

func (h *BaseHandler) ParamError(ctx *gin.Context, err error) {
	ParamError(ctx, err)
}

func (h *BaseHandler) Custom(ctx *gin.Context, httpCode int, data *Response) {
	Custom(ctx, httpCode, data)
}
