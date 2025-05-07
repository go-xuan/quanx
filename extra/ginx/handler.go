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

// BindRouter 绑定路由
// handler必须实现 Handler 接口，按需实现 GetHandler, PostHandler ... 等接口
func BindRouter(router *gin.RouterGroup, handler Handler, middleware ...gin.HandlerFunc) {
	// 设置路由和中间件
	group := router.Group(handler.RelativePath(), middleware...)
	if impl, ok := handler.(GetHandler); ok {
		group.GET("", impl.Get)
	}
	if impl, ok := handler.(PostHandler); ok {
		group.POST("", impl.Post)
	}
	if impl, ok := handler.(PutHandler); ok {
		group.PUT("", impl.Put)
	}
	if impl, ok := handler.(DeleteHandler); ok {
		group.DELETE("", impl.Delete)
	}
	if impl, ok := handler.(PatchHandler); ok {
		group.PATCH("", impl.Patch)
	}
	if impl, ok := handler.(HeadHandler); ok {
		group.HEAD("", impl.Head)
	}
	if impl, ok := handler.(OptionsHandler); ok {
		group.OPTIONS("", impl.Options)
	}
	if impl, ok := handler.(AnyHandler); ok {
		group.Any("", impl.Any)
	}
}
