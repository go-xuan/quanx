package mongox

import (
	"testing"

	"github.com/go-xuan/quanx/configx"
)

func TestMongo(t *testing.T) {
	if err := configx.ReadAndExecute(&Config{}, configx.FromFile); err != nil {
		panic(err)
	}
}
