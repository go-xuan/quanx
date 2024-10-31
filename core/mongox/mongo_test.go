package mongox

import "testing"

func TestMongo(t *testing.T) {
	if err := NewConfigurator(&Config{
		URI:      "",
		Username: "",
		Password: "",
		Database: "",
	}).Execute(); err != nil {
		t.Error(err)
	}
}
