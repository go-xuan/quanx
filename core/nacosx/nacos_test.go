package nacosx

import "testing"

func TestNacos(t *testing.T) {
	if err := NewConfigurator(&Config{
		Address:   "localhost:8848",
		Username:  "nacos",
		Password:  "nacos",
		NameSpace: "",
		Mode:      0,
	}); err != nil {
		t.Error(err)
	}
}
