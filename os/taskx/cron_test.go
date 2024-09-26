package taskx

import (
	"fmt"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	c := Corn()

	var count1, count2 = 0, 0
	if err := c.Add("task1", "@every 3s", func() {
		count1++
		fmt.Printf("第%d次执行task1\n", count1)
	}); err != nil {
		fmt.Println(err)
	}

	if err := c.Add("task2", "@every 5s", func() {
		count2++
		fmt.Printf("第%d次执行task2\n", count2)
	}); err != nil {
		fmt.Println(err)
	}

	c.Start()
	time.Sleep(time.Minute)
	c.Stop()
}
