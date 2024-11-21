package timex

import "time"

func NewDate(v ...time.Time) Date {
	var x = Date{notnull: true}
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
	if value, err := time.ParseInLocation(`"2006-01-02"`, string(bytes), time.Local); err != nil {
		return err
	} else {
		x.value = value
		x.notnull = true
	}
	return nil
}

func (x *Date) MarshalJSON() ([]byte, error) {
	if x.notnull {
		var bytes []byte
		bytes = append(bytes, 34)
		bytes = x.value.AppendFormat(bytes, DateFmt)
		bytes = append(bytes, 34)
		return bytes, nil
	} else {
		return []byte("null"), nil
	}
}

func (x *Date) Value(def ...time.Time) time.Time {
	if x.notnull {
		return x.value
	} else if len(def) > 0 {
		y, m, d := def[0].Date()
		return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	}
	return time.Time{}
}

func (x *Date) NotNull() bool {
	return x.notnull
}
