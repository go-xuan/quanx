package timex

import (
	"strconv"
	"time"
)

func NewTime(v ...time.Time) *Time {
	var x = &Time{notnull: true}
	if len(v) > 0 {
		x.value = v[0]
	} else {
		x.value = time.Now()
	}
	return x
}

type Time struct {
	value   time.Time
	notnull bool
}

func (x *Time) UnmarshalJSON(bytes []byte) error {
	if l := len(bytes); l > 1 && string(bytes) != "null" {
		if bytes[0] == 34 && bytes[l-1] == 34 {
			bytes = bytes[1 : l-1] // 带引号则去掉引号
		}
		str := string(bytes)
		if t, err := time.ParseInLocation(`2006-01-02 15:04:05`, str, time.Local); err == nil {
			x.value = t
			x.notnull = true
			return nil
		}
		if unix, err := strconv.ParseInt(str, 10, 64); err == nil {
			if unix > 1e12 {
				x.value = time.UnixMilli(unix)
			} else {
				x.value = time.Unix(unix, 0)
			}
			x.notnull = true
			return nil
		}
	}
	x.notnull = false
	return nil
}

func (x *Time) MarshalJSON() ([]byte, error) {
	if x.Valid() {
		var bytes []byte
		bytes = append(bytes, 34)
		bytes = x.value.AppendFormat(bytes, TimeFmt)
		bytes = append(bytes, 34)
		return bytes, nil
	}
	return []byte("null"), nil
}

func (x *Time) Value(def ...time.Time) time.Time {
	if x.Valid() {
		return x.value
	} else if len(def) > 0 {
		return def[0]
	}
	return time.Time{}
}

func (x *Time) Valid() bool {
	return x != nil && x.notnull
}
