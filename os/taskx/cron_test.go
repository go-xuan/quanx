package taskx

import (
	"fmt"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	c := Corn()
	if err := c.Add("task1", "@every 3s", func() {
		fmt.Println("task1")
	}); err != nil {
		fmt.Println(err)
	}
	c.Start()
	time.Sleep(10 * time.Second)
	c.Stop()
}
