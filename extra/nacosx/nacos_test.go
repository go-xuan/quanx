package nacosx

import (
	"testing"

	"github.com/go-xuan/quanx/extra/configx"
)

func TestNacos(t *testing.T) {
	if err := configx.ReadAndExecute(&Config{
		Username: "nacos",
		Password: "nacos",
	}, configx.FromDefault); err != nil {
		panic(err)
	}
}
