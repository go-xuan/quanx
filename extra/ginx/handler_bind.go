package ginx

import (
	"github.com/gin-gonic/gin"
)

func BindRouter(router *gin.RouterGroup, handler Handler, handlers ...gin.HandlerFunc) {
	// 设置路由和中间件
	group := router.Group(handler.RelativePath(), handlers...)
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
