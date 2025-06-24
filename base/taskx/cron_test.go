package taskx

import (
	"context"
	"fmt"
	"testing"
)

func job1(context.Context) {
	fmt.Println("执行定时任务 ==> job1")
}

func job2(context.Context) {
	fmt.Println("执行定时任务 ==> job2")
}

func job3(context.Context) {
	fmt.Println("执行定时任务 ==> job3")
}

func job4(context.Context) {
	fmt.Println("执行定时任务 ==> job4")
}

func TestCron(t *testing.T) {
	// 初始化
	scheduler := Cron(
		DurationWrap,
	)

	scheduler.Add("job1", "@every 5s", job1)
	scheduler.Add("job2", "@every 2s", job2)
	scheduler.Add("job3", "@daily", job3)
	scheduler.Add("job4", "0 */1 * * * ?", job4)

	scheduler.Start()

	// 定时任务信息
	for i, job := range scheduler.All() {
		fmt.Println(i, job.Info())
	}
	select {}
}
