package nacosx

import (
	"testing"

	"github.com/go-xuan/quanx/extra/configx"
)

func TestNacos(t *testing.T) {
	if err := configx.ReadAndExecute(&Config{}, configx.FromFile, ""); err != nil {
		panic(err)
	}
}
