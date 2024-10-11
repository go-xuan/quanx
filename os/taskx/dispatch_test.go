package taskx

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-xuan/quanx/utils/randx"
)

func TestInBatches(t *testing.T) {
	if err := InBatches(100, 10, func(x int, y int) error {
		fmt.Println(x, "==>", y)
		time.Sleep(1 * time.Second)
		return nil
	}); err != nil {
		t.Error(err)
	}
}

func TestRetry(t *testing.T) {
	if err := Retry(10, 0, time.Second*1, func() error {
		// todo 执行处理逻辑

		// 五分之一概率随机成功
		if randx.IntRange(1, 5) == 1 {
			return nil
		}
		return errors.New("error")
	}); err != nil {
		t.Error(err)
	}
}
