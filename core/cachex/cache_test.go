package cachex

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/core/redisx"
)

func TestHandler(t *testing.T) {
	// 先初始化redis
	if err := configx.Execute(&redisx.Config{
		Source:   "default",
		Enable:   true,
		Host:     "localhost",
		Port:     6379,
		Username: "",
		Password: "Init@1234",
		Database: 1,
		PoolSize: 15,
	}); err != nil {
		panic(err)
	}

	// 再初始化缓存客户端
	if err := configx.Execute(&Config{
		Type:    "redis",
		Source:  "default",
		Prefix:  "test",
		Marshal: "json",
	}); err != nil {
		panic(err)
	}

	ctx := context.TODO()
	Client().Set(ctx, "test_1", "111", -1)

	var value string
	Client().Get(ctx, "test_1", &value)
	fmt.Println(value)
}
