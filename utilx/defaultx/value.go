package defaultx

import (
	"time"
)

func Int(value, def int) int {
	if value == 0 {
		return def
	}
	return value
}

func Int64(value, def int64) int64 {
	if value == 0 {
		return def
	}
	return value
}

func Int32(value, def int32) int32 {
	if value == 0 {
		return def
	}
	return value
}

func String(value, def string) string {
	if value == "" {
		return def
	}
	return value
}

func Time(value, def time.Time) time.Time {
	if value.IsZero() {
		return def
	}
	return value
}

func Any(value, def any) any {
	if value == nil {
		return def
	}
	return value
}
