package logx

import (
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
)

func TestLog(t *testing.T) {
	log.WithField("test", "test").Info("test")
	if err := configx.ConfiguratorReadAndExecute(&Config{}); err != nil {
		panic(err)
	}
	log.WithField("test", "test").Info("test")
}
