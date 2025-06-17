package timex

import (
	"time"
)

func NewDate(v ...time.Time) *Date {
	var x = &Date{notnull: true}
	if len(v) > 0 {
		y, m, d := v[0].Date()
		x.value = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	} else {
		y, m, d := time.Now().Date()
		x.value = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	}
	return x
}

type Date struct {
	value   time.Time
	notnull bool
}

func (x *Date) UnmarshalJSON(bytes []byte) error {
	if l := len(bytes); l > 1 && string(bytes) != "null" {
		if bytes[0] == 34 && bytes[l-1] == 34 {
			bytes = bytes[1 : l-1] // 带引号则去掉引号
		}
		str := string(bytes)
		if t, err := time.ParseInLocation(`2006-01-02`, str, time.Local); err == nil {
			x.value = t
			x.notnull = true
			return nil
		}
	}
	x.notnull = false
	return nil
}

func (x *Date) MarshalJSON() ([]byte, error) {
	if x.Valid() {
		var bytes []byte
		bytes = append(bytes, 34)
		bytes = x.value.AppendFormat(bytes, DateFmt)
		bytes = append(bytes, 34)
		return bytes, nil
	}
	return []byte("null"), nil
}

func (x *Date) Value(def ...time.Time) time.Time {
	if x.Valid() {
		return x.value
	} else if len(def) > 0 {
		y, m, d := def[0].Date()
		return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	}
	return time.Time{}
}

func (x *Date) Valid() bool {
	return x != nil && x.notnull
}
