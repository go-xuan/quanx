package elasticx

import (
	"testing"

	"github.com/go-xuan/quanx/configx"
)

func TestElastic(t *testing.T) {
	if err := configx.LoadConfigurator(&Config{}); err != nil {
		panic(err)
	}
}
