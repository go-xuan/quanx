package taskx

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := Queue()
	queue.Add("task1", func() error { return nil })                // 正常插入 task1
	queue.Add("task1", func() error { return nil })                // 默认尾插 task1->task2
	queue.AddAfter("task3", "task1", func() error { return nil })  // 插队到task1后面 task1->task3->task2
	queue.Add("task4", func() error { return nil })                // 默认尾插 task1->task3->task2->task4
	queue.AddBefore("task5", "task4", func() error { return nil }) // 插队到task4前面 task1->task3->task2->task5->task4
	queue.AddTail("task6", func() error { return nil })            // 插队到末位 task1->task3->task2->task5->task4->task6
	queue.AddHead("task7", func() error { return nil })            // 插队到首位 task7->task1->task3->task2->task5->task4->task6
	queue.Remove("task2")                                          // 删除task2 task7->task1->task3->task5->task4->task6
	fmt.Println(queue.Names())
	queue.Execute()
}
