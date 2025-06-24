package taskx

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := Queue()
	queue.Add("job1", func() error { return nil })                 // 正常插入 job1
	queue.Add("job1", func() error { return nil })                 // 默认尾插 job1->task2
	queue.AddAfter("job1", "task3", func() error { return nil })   // 插队到task1后面 job1->task3->task2
	queue.Add("task4", func() error { return nil })                // 默认尾插 job1->task3->task2->task4
	queue.AddBefore("task4", "task5", func() error { return nil }) // 插队到task4前面 job1->task3->task2->task5->task4
	queue.AddTail("task6", func() error { return nil })            // 插队到末位 job1->task3->task2->task5->task4->task6
	queue.AddHead("task7", func() error { return nil })            // 插队到首位 task7->job1->task3->task2->task5->task4->task6
	queue.Remove("task2")                                          // 删除task2 task7->job1->task3->task5->task4->task6
	fmt.Println(queue.Names())
	queue.Execute()
}
