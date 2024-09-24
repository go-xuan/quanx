package taskx

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	q := Queue()

	q.Add("task1", func() {
		fmt.Println("task1")
	})

	q.Add("task2", func() {
		fmt.Println("task2")
	})

	q.AddAfter("task3", func() {
		fmt.Println("task3")
	}, "task1")

	q.Add("task4", func() {
		fmt.Println("task4")
	})

	q.AddBefore("task5", func() {
		fmt.Println("task5")
	}, "task4")

	q.AddTail("task6", func() {
		fmt.Println("task6")
	})

	q.AddHead("task7", func() {
		fmt.Println("task7")
	})

	q.Execute()
}
