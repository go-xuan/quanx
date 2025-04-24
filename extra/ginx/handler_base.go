package ginx

import (
	"github.com/gin-gonic/gin"
	
	"github.com/go-xuan/quanx/base/respx"
)

type BaseHandler struct {
	Path string
}

func (h *BaseHandler) RelativePath() string {
	return h.Path
}

func (h *BaseHandler) Success(ctx *gin.Context, data any) {
	respx.Success(ctx, data)
}

func (h *BaseHandler) Error(ctx *gin.Context, data any) {
	respx.Error(ctx, data)
}

func (h *BaseHandler) ParamError(ctx *gin.Context, err error) {
	respx.ParamError(ctx, err)
}
