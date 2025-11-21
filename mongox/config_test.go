package mongox

import (
	"testing"

	"github.com/go-xuan/quanx/configx"
)

func TestMongo(t *testing.T) {
	if err := configx.LoadConfigurator(&Config{}); err != nil {
		panic(err)
	}
}
