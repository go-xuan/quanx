package mongox

import (
	"testing"
)

func TestMongo(t *testing.T) {
	if err := (&Config{
		URI:      "mongodb://localhost:27017",
		Username: "",
		Password: "",
		Database: "",
	}).Execute(); err != nil {
		t.Error(err)
	}
}
