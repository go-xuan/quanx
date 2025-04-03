package elasticx

import (
	"testing"

	"github.com/go-xuan/quanx/extra/configx"
)

func TestElastic(t *testing.T) {
	if err := configx.ReadAndExecute(&Config{}, configx.FromDefault); err != nil {
		panic(err)
	}
}
