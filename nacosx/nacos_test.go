package nacosx

import (
	"testing"
	
	"github.com/go-xuan/quanx/configx"
)

func TestNacos(t *testing.T) {
	if err := configx.ReadAndExecute(&Config{}, configx.FromFile, ""); err != nil {
		panic(err)
	}
}
