package ossx

import (
	"testing"

	"github.com/go-xuan/quanx/configx"
)

func TestOss(t *testing.T) {
	if err := configx.LoadConfigurator(&Config{}); err != nil {
		panic(err)
	}
}
