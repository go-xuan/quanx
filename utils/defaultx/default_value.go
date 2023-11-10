package defaultx

import (
	"time"

	"github.com/go-xuan/quanx/utils/timex"
)

func Int(target, def int) int {
	if target == 0 {
		return def
	}
	return target
}

func Int64(target, def int64) int64 {
	if target == 0 {
		return def
	}
	return target
}

func Int32(target, def int32) int32 {
	if target == 0 {
		return def
	}
	return target
}

func String(target, def string) string {
	if target == "" {
		return def
	}
	return target
}

func Time(target string, def time.Time) time.Time {
	if target == "" {
		return def
	}
	return timex.ToTime(target)
}

func Any(target, def any) any {
	if target == nil {
		return def
	}
	return target
}
