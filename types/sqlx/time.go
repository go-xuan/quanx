package sqlx

import (
	"time"

	"github.com/go-xuan/quanx/types/timex"
)

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	var location time.Time
	if location, err = time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local); err != nil {
		return
	}
	*t = Time(location)
	return
}

func (t *Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timex.TimeFmt)+2)
	b = append(b, '"')
	b = time.Time(*t).AppendFormat(b, timex.TimeFmt)
	b = append(b, '"')
	return b, nil
}

func (t *Time) String() string {
	return time.Time(*t).Format(timex.TimeFmt)
}

type Date time.Time

func (t *Date) UnmarshalJSON(data []byte) (err error) {
	var location time.Time
	if location, err = time.ParseInLocation(`"2006-01-02"`, string(data), time.Local); err != nil {
		return
	}
	*t = Date(location)
	return
}

func (t *Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timex.DateFmt)+2)
	b = append(b, '"')
	b = time.Time(*t).AppendFormat(b, timex.DateFmt)
	b = append(b, '"')
	return b, nil
}

func (t *Date) String() string {
	return time.Time(*t).Format(timex.DateFmt)
}
