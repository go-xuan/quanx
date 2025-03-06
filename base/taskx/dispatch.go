package taskx

import (
	"fmt"
	"time"

	"github.com/go-xuan/quanx/base/errorx"
)

// ExecWithBatches 分批次执行
func ExecWithBatches(total, limit int, f func(int, int) error) error {
	var start, end int
	for start < total {
		if end = start + limit; end > total {
			end = total
		}
		if err := f(start, end); err != nil {
			return errorx.Wrap(err, fmt.Sprintf("execute failed between %d and %d", start, end))
		}
		start = end
	}
	return nil
}

// Retry 重试执行
func Retry(times, index int, interval time.Duration, f func() error) error {
	if times > 0 {
		time.Sleep(interval)
		index++
		if err := f(); err != nil {
			fmt.Printf("execute failed after %d times retries\n", index)
			return Retry(times-1, index, interval, f)
		} else {
			fmt.Printf("execute success after %d times retries\n", index)
			return nil
		}
	}
	return errorx.New("retry failed")
}
