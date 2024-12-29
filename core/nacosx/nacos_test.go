package nacosx

import (
	"github.com/go-xuan/quanx/core/configx"
	"testing"
)

func TestNacos(t *testing.T) {
	if err := configx.Execute(&Config{
		Address:   "localhost:8848",
		Username:  "nacos",
		Password:  "nacos",
		NameSpace: "",
		Mode:      0,
	}); err != nil {
		t.Error(err)
	}
}
