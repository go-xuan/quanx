package dbx

import (
	"testing"

	"github.com/go-xuan/quanx/configx"
)

func TestDatabase(t *testing.T) {
	if err := configx.LoadConfigurator(&Config{}); err != nil {
		panic(err)
	}
}
