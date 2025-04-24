package taskx

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := Queue()
	queue.Add(func() error { return nil }, "task1")                // 正常插入 task1
	queue.Add(func() error { return nil }, "task2")                // 默认尾插 task1->task2
	queue.AddAfter(func() error { return nil }, "task3", "task1")  // 插队到task1后面 task1->task3->task2
	queue.Add(func() error { return nil }, "task4")                // 默认尾插 task1->task3->task2->task4
	queue.AddBefore(func() error { return nil }, "task5", "task4") // 插队到task4前面 task1->task3->task2->task5->task4
	queue.AddTail(func() error { return nil }, "task6")            // 插队到末位 task1->task3->task2->task5->task4->task6
	queue.AddHead(func() error { return nil }, "task7")            // 插队到首位 task7->task1->task3->task2->task5->task4->task6
	queue.Remove("task2")                                          // 删除task2 task7->task1->task3->task5->task4->task6
	fmt.Println(queue.Names())
	queue.Execute()
}
