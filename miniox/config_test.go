package miniox

import (
	"testing"

	"github.com/go-xuan/quanx/configx"
)

func TestMinio(t *testing.T) {
	if err := configx.ConfiguratorReadAndExecute(&Config{}); err != nil {
		panic(err)
	}
}
