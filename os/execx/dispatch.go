package execx

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
				return err
			}
			start += limit
		}
	}
	return nil
}

// Retry 重试执行
func Retry(times, count int, interval time.Duration, f func() error) error {
	if times > 0 {
		time.Sleep(interval)
		count++
		if err := f(); err != nil {
			fmt.Printf("重试第%d次执行失败！\n", count)
			return Retry(times-1, count, interval, f)
		} else {
			fmt.Printf("重试第%d次执行成功！\n", count)
			return nil
		}
	}
	return errorx.New("重试失败")
}
