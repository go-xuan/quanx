package redisx

import (
	"testing"

	"github.com/go-xuan/quanx/configx"
)

func TestRedis(t *testing.T) {
	if err := configx.ConfiguratorReadAndExecute(&Config{}); err != nil {
		panic(err)
	}

	if err := CopyDatabase("default", "target", 1); err != nil {
		panic(err)
	}
}
