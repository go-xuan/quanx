package ginx

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	RelativePath() string
}

type GetHandler interface {
	Get(ctx *gin.Context)
}

type PostHandler interface {
	Post(ctx *gin.Context)
}

type PutHandler interface {
	Put(ctx *gin.Context)
}

type DeleteHandler interface {
	Delete(ctx *gin.Context)
}

type PatchHandler interface {
	Patch(ctx *gin.Context)
}

type HeadHandler interface {
	Head(ctx *gin.Context)
}

type OptionsHandler interface {
	Options(ctx *gin.Context)
}

type AnyHandler interface {
	Any(ctx *gin.Context)
}
