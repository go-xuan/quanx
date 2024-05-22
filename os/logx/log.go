package logx

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func UseTime(name string, f func()) {
	start := time.Now()
	log.Info(name, "start")
	f()
	log.Infof("%s finish use %dms", name, time.Since(start).Milliseconds())
}
