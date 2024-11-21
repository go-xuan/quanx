package timex

import (
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

func (v *Time) UnmarshalJSON(bytes []byte) error {
	if value, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(bytes), time.Local); err != nil {
		return err
	} else {
		v.value = value
		v.notnull = true
	}
	return nil
}

func (v *Time) MarshalJSON() ([]byte, error) {
	if v.notnull {
		var bytes []byte
		bytes = append(bytes, 34)
		bytes = v.value.AppendFormat(bytes, TimeFmt)
		bytes = append(bytes, 34)
		return bytes, nil
	} else {
		return []byte("null"), nil
	}
}

func (v *Time) Value() time.Time {
	return v.value
}

func (v *Time) NotNull() bool {
	return v.notnull
}
