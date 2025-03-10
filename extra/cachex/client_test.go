package cachex

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-xuan/quanx/extra/redisx"
)

func TestHandler(t *testing.T) {
	// 先初始化redis
	if err := (&redisx.Config{
		Source:   "default",
		Enable:   true,
		Host:     "localhost",
		Port:     6379,
		Username: "",
		Password: "Init@1234",
		Database: 1,
		PoolSize: 15,
	}).Execute(); err != nil {
		fmt.Println(err)
		return
	}

	// 再初始化缓存客户端
	if err := (&Config{
		Type:   "redis",
		Source: "default",
		Prefix: "test_",
	}).Execute(); err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.TODO()
	if err := GetClient().Set(ctx, "test_1", "111", -1); err != nil {
		fmt.Println(err)
		return
	}

	var value string
	GetClient().Get(ctx, "test_1", &value)
	fmt.Println(value)
}
