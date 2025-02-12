package mongox

import (
	"github.com/go-xuan/quanx/core/configx"
	"testing"
)

func TestMongo(t *testing.T) {
	if err := configx.Execute(&Config{
		URI:      "mongodb://localhost:27017",
		Username: "",
		Password: "",
		Database: "",
	}); err != nil {
		t.Error(err)
	}
}
