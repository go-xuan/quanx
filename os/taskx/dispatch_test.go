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
	if err := Retry(5, 0, time.Second, func() error {
		if randx.Bool() {
			// 随机成功
			return nil
		}
		return errors.New("error")
	}); err != nil {
		t.Error(err)
	}
}
