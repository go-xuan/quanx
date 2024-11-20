package timex

import (
	"time"
)

func NewTime(t ...time.Time) Time {
	var ti = Time{notnull: true}
	if len(t) > 0 {
		ti.value = t[0]
	} else {
		ti.value = time.Now()
	}
	return ti
}

func NewDate(t ...time.Time) Date {
	var date = Date{notnull: true}
	if len(t) > 0 {
		y, m, d := t[0].Date()
		date.value = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	} else {
		y, m, d := time.Now().Date()
		date.value = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	}
	return date
}

type Time struct {
	value   time.Time
	notnull bool
}

func (t *Time) UnmarshalJSON(bytes []byte) error {
	if value, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(bytes), time.Local); err != nil {
		return err
	} else {
		t.value = value
		t.notnull = true
	}
	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if t.notnull {
		var bytes []byte
		bytes = append(bytes, 34)
		bytes = t.value.AppendFormat(bytes, TimeFmt)
		bytes = append(bytes, 34)
		return bytes, nil
	} else {
		return []byte("null"), nil
	}
}

func (t *Time) Value() time.Time {
	return t.value
}

func (t *Time) NotNull() bool {
	return t.notnull
}

type Date struct {
	value   time.Time
	notnull bool
}

func (d *Date) UnmarshalJSON(bytes []byte) error {
	if value, err := time.ParseInLocation(`"2006-01-02"`, string(bytes), time.Local); err != nil {
		return err
	} else {
		d.value = value
		d.notnull = true
	}
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	if d.notnull {
		var bytes []byte
		bytes = append(bytes, 34)
		bytes = d.value.AppendFormat(bytes, DateFmt)
		bytes = append(bytes, 34)
		return bytes, nil
	} else {
		return []byte("null"), nil
	}
}

func (d *Date) Value() time.Time {
	return d.value
}

func (d *Date) NotNull() bool {
	return d.notnull
}
