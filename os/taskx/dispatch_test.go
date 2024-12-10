package taskx

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-xuan/quanx/utils/randx"
)

func TestInBatches(t *testing.T) {
	var total = 100
	var s []int
	for i := 0; i < total; i++ {
		s = append(s, i)
	}
	if err := ExecWithBatches(total, 10, func(start, end int) error {
		fmt.Printf("%d ==> %d :%v \n", start, end, s[start:end])
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
