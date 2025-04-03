package cachex

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-xuan/quanx/extra/configx"
)

func TestCache(t *testing.T) {
	if err := configx.ReadAndExecute(&Config{}, configx.FromDefault); err != nil {
		panic(err)
	}

	ctx := context.TODO()
	if err := GetClient().Set(ctx, "test_1", "111", -1); err != nil {
		fmt.Println(err)
		return
	}

	var value string
	fmt.Println("value:", value)
	GetClient().Get(ctx, "test_1", &value)
	fmt.Println("value:", value)
}
