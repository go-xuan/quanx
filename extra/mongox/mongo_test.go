package mongox

import (
	"testing"

	"github.com/go-xuan/quanx/extra/configx"
)

func TestMongo(t *testing.T) {
	if err := configx.ReadAndExecute(&Config{
		Enable:   true,
		URI:      "mongodb://localhost:27017",
		Database: "test",
	}, configx.FromDefault); err != nil {
		panic(err)
	}
}
