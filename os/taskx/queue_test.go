package taskx

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := Queue()

	queue.Add(func() {
		fmt.Println("task1")
	}, "task1")

	queue.Add(func() {
		fmt.Println("task2")
	}, "task2")

	// 插队到目标位置后一位
	queue.AddAfter(func() {
		fmt.Println("task3")
	}, "task3", "task1")

	queue.Add(func() {
		fmt.Println("task4")
	}, "task4")

	// 插队到目标位置前一位
	queue.AddBefore(func() {
		fmt.Println("task5")
	}, "task5", "task4")

	// 插队到末位
	queue.AddTail(func() {
		fmt.Println("task6")
	}, "task6")

	// 插队到首位
	queue.AddHead(func() {
		fmt.Println("task7")
	}, "task7")

	queue.Execute()
}
