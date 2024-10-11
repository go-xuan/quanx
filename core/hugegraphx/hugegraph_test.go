package hugegraphx

import "testing"

func TestHugegraph(t *testing.T) {
	if err := NewConfigurator(&Config{
		Host:  "localhost",
		Port:  8882,
		Graph: "hugegraph",
	}).Execute(); err != nil {
		t.Error(err)
	}
}
