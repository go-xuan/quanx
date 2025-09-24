package gormx

import (
	"testing"

	"github.com/go-xuan/quanx/configx"
)

func TestDatabase(t *testing.T) {
	if err := configx.ConfiguratorReadAndExecute(&Config{}); err != nil {
		panic(err)
	}
}
