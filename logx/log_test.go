package logx

import (
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
)

func TestLog(t *testing.T) {
	log.WithField("test", "test").Info("test")
	if err := configx.ReadAndExecute(&Config{}, configx.FromFile); err != nil {
		panic(err)
	}
	log.WithField("test", "test").Info("test")
}
