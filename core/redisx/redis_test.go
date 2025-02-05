package redisx

import (
	"context"
	"fmt"
	"github.com/go-xuan/quanx/core/configx"
	"testing"
)

func TestHandler(t *testing.T) {
	// 先初始化缓存
	if err := configx.Execute(&Config{
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
	ctx := context.TODO()
	GetClient().Set(ctx, "test_1", "111", -1)
	value := GetClient().Get(ctx, "test_1").Val()
	fmt.Println(value)
}
