package logx

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func Ctx(ctx context.Context) *log.Entry {

	return log.WithContext(ctx)
}
