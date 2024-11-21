package timex

import (
	"strconv"
	"time"
)

func NewTime(v ...time.Time) Time {
	var x = Time{notnull: true}
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
	if l := len(bytes); l >= 0 {
		if l > 1 && bytes[0] == 34 && bytes[l-1] == 34 {
			if value, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(bytes), time.Local); err != nil {
				return err
			} else {
				x.value = value
				x.notnull = true
			}
		} else { // 兼容时间戳
			if value, err := strconv.ParseInt(string(bytes), 10, 64); err != nil {
				return err
			} else {
				x.value = time.Unix(value, 0)
				x.notnull = true
			}
		}
	} else {
		x.notnull = false
	}
	return nil
}

func (x *Time) MarshalJSON() ([]byte, error) {
	if x.notnull {
		var bytes []byte
		bytes = append(bytes, 34)
		bytes = x.value.AppendFormat(bytes, TimeFmt)
		bytes = append(bytes, 34)
		return bytes, nil
	} else {
		return []byte("null"), nil
	}
}

func (x *Time) Value(def ...time.Time) time.Time {
	if x.notnull {
		return x.value
	} else if len(def) > 0 {
		return def[0]
	}
	return time.Time{}
}

func (x *Time) NotNull() bool {
	return x.notnull
}
