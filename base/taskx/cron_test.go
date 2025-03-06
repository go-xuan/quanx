package taskx

import (
	"fmt"
	"testing"
)

func task1() {
	fmt.Println("执行定时任务 ==> task1")
}

func task2() {
	fmt.Println("执行定时任务 ==> task2")
}

func task3() {
	fmt.Println("执行定时任务 ==> task3")
}

func task4() {
	fmt.Println("执行定时任务 ==> task4")
}

func TestCron(t *testing.T) {
	// 初始化
	scheduler := Corn(
		DurationWarp,
	)

	//scheduler.Add("task1", "@every 5s", task1)
	//scheduler.Add("task2", "@every 2s", task2)
	//scheduler.Add("task3", "@daily", task3)
	//scheduler.Add("task4", "0 */1 * * * ?", task4)

	scheduler.Start()

	// 定时任务信息
	for i, task := range scheduler.All() {
		fmt.Println(i, task.Info())
	}
	select {}
}
