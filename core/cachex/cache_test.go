package cachex

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-xuan/quanx/core/redisx"
)

func TestHandler(t *testing.T) {
	// 先初始化redis
	if err := redisx.NewConfigurator(&redisx.Config{
		Source:   "default",
		Enable:   true,
		Host:     "localhost",
		Port:     6379,
		Username: "",
		Password: "Init@1234",
		Database: 1,
		PoolSize: 15,
	}).Execute(); err != nil {
		panic(err)
	}

	// 再初始化缓存客户端
	if err := NewConfigurator(&Config{
		Type:    "redis",
		Source:  "default",
		Prefix:  "test",
		Marshal: "json",
	}).Execute(); err != nil {
		panic(err)
	}

	ctx := context.TODO()
	Client().Set(ctx, "test_1", "111", -1)

	var value string
	Client().Get(ctx, "test_1", &value)
	fmt.Println(value)
}
