package taskx

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := Queue()

	queue.Add("task1", func() {
		fmt.Println("task1")
	})

	queue.Add("task2", func() {
		fmt.Println("task2")
	})

	// 插队到目标位置后一位
	queue.AddAfter("task3", func() {
		fmt.Println("task3")
	}, "task1")

	queue.Add("task4", func() {
		fmt.Println("task4")
	})

	// 插队到目标位置前一位
	queue.AddBefore("task5", func() {
		fmt.Println("task5")
	}, "task4")

	// 插队到末位
	queue.AddTail("task6", func() {
		fmt.Println("task6")
	})

	// 插队到首位
	queue.AddHead("task7", func() {
		fmt.Println("task7")
	})

	queue.Execute()
}
