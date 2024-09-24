package taskx

import (
	"fmt"
	"time"

	"github.com/go-xuan/quanx/os/errorx"
)

// InBatches 分批次执行
func InBatches(total, limit int, f func(int, int) error) error {
	var start = 0
	if start < total {
		for start < total {
			if start+limit > total {
				limit = total - start
			}
			if err := f(start, start+limit); err != nil {
				return errorx.Wrap(err, "execute failed")
			}
			start += limit
		}
	}
	return nil
}

// Retry 重试执行
func Retry(size, index int, interval time.Duration, f func() error) error {
	if size > 0 {
		time.Sleep(interval)
		index++
		if err := f(); err != nil {
			fmt.Printf("execute failed after %d times retery\n", index)
			return Retry(size-1, index, interval, f)
		} else {
			fmt.Printf("execute success after %d times retery\n", index)
			return nil
		}
	}
	return errorx.New("retry failed")
}
