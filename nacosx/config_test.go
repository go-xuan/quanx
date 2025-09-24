package nacosx

import (
	"testing"

	"github.com/go-xuan/quanx/configx"
)

func TestNacos(t *testing.T) {
	if err := configx.ConfiguratorReadAndExecute(&Config{}); err != nil {
		panic(err)
	}
}
