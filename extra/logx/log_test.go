package logx

import (
	"testing"

	log "github.com/sirupsen/logrus"
	
	"github.com/go-xuan/quanx/extra/configx"
)

func TestLog(t *testing.T) {
	log.WithField("test", "test").Info("test")
	if err := configx.ReadAndExecute(&Config{}, configx.FromDefault); err != nil {
		panic(err)
	}
	log.WithField("test", "test").Info("test")
}
