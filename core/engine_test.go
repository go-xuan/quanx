package core

import (
	"fmt"
	"testing"
)

func TestEngineRun(t *testing.T) {
	var newEngine = GetEngine(NonGin)

	// 添加初始化方法
	newEngine.AddCustomFunc(
		func() {
			fmt.Println("初始化加载方法1")
		},
		func() {
			fmt.Println("初始化加载方法2")
		},
	)

	// 服务启动
	newEngine.RUN()
}
