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

func (v *Date) UnmarshalJSON(bytes []byte) error {
	if value, err := time.ParseInLocation(`"2006-01-02"`, string(bytes), time.Local); err != nil {
		return err
	} else {
		v.value = value
		v.notnull = true
	}
	return nil
}

func (v *Date) MarshalJSON() ([]byte, error) {
	if v.notnull {
		var bytes []byte
		bytes = append(bytes, 34)
		bytes = v.value.AppendFormat(bytes, DateFmt)
		bytes = append(bytes, 34)
		return bytes, nil
	} else {
		return []byte("null"), nil
	}
}

func (v *Date) Value() time.Time {
	return v.value
}

func (v *Date) NotNull() bool {
	return v.notnull
}
