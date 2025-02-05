package logx

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/core/ginx"
)

func Trace(ctx *gin.Context) *log.Entry {
	return log.WithField("traceId", ginx.TraceId(ctx))
}
