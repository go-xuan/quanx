package logx

import (
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func Ctx(ctx *gin.Context) *log.Entry {
	entry := log.WithContext(ctx)
	var fields = entry.Data
	for k, v := range ctx.Keys {
		fields[k] = v
	}
	entry.Data = fields
	return entry
}
