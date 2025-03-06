package elasticx

import (
	"testing"
)

func TestElastic(t *testing.T) {
	if err := (&Config{
		Host: "localhost",
		Port: 9200,
	}).Execute(); err != nil {
		t.Error(err)
	}
}
