package logx

import (
	"context"
	
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/core/ginx"
)

func GinCtx(ctx *gin.Context) *log.Entry {
	return ginx.Log(ctx)
}

func Ctx(ctx context.Context) *log.Entry {
	return log.WithContext(ctx)
}
