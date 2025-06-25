package redisx

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-xuan/quanx/configx"
)

func TestRedis(t *testing.T) {
	if err := configx.ReadAndExecute(&Config{}, configx.FromTag); err != nil {
		panic(err)
	}

	if err := CopyDatabase("default", "target", 1); err != nil {
		panic(err)
	}

	if ok, err := Ping(context.TODO(), "target"); err != nil {
		panic(err)
	} else {
		fmt.Println(ok)
	}
}
