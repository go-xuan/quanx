package mongox

import (
	"testing"

	"github.com/go-xuan/quanx/configx"
)

func TestMongo(t *testing.T) {
	if err := configx.ConfiguratorReadAndExecute(&Config{}); err != nil {
		panic(err)
	}
}
