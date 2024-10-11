package mongox

import "testing"

func TestMongo(t *testing.T) {
	if err := NewConfigurator(&Config{
		Host:     "",
		Port:     0,
		Username: "",
		Password: "",
		Database: "",
	}).Execute(); err != nil {
		t.Error(err)
	}
}
