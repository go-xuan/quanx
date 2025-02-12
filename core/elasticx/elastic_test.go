package elasticx

import (
	"testing"

	"github.com/go-xuan/quanx/core/configx"
)

func TestElastic(t *testing.T) {
	if err := configx.Execute(&Config{
		Host: "localhost",
		Port: 9200,
	}); err != nil {
		t.Error(err)
	}
}
