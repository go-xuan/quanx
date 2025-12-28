package redisx

import (
	"testing"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/constx"
)

func TestRedis(t *testing.T) {
	if err := configx.LoadConfigurator(&Config{}); err != nil {
		panic(err)
	}

	if err := CopyDatabase(constx.DefaultSource, "target", 1); err != nil {
		panic(err)
	}
}
