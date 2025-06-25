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

func (h *BaseHandler) Error(ctx *gin.Context, data any) {
	Error(ctx, data)
}

func (h *BaseHandler) ParamError(ctx *gin.Context, err error) {
	ParamError(ctx, err)
}
