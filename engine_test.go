package quanx

import (
	"fmt"
	"testing"
)

func TestEngineRun(t *testing.T) {
	var e = NewEngine(NonGin)

	// 添加初始化方法
	e.AddCustomFunc(
		func() {
			fmt.Println("初始化加载方法1")
		},
		func() {
			fmt.Println("初始化加载方法2")
		},
	)

	e.AddQueueTask("task1", func() {
		fmt.Println("队列任务1")
	})

	e.AddQueueTask("task1", func() {
		fmt.Println("队列任务2")
	})

	e.AddQueueTask("task1", func() {
		fmt.Println("队列任务3")
	})

	// 服务启动
	e.RUN()
}
