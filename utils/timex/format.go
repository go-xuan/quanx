package timex

import (
	"time"

	"github.com/quanxiaoxuan/quanx/common/constx"
)

type Time time.Time
type Date time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	tt, err := time.ParseInLocation(`"`+constx.TimeFmt+`"`, string(data), time.Local)
	if err != nil {
		return err
	}
	*t = Time(tt)
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(constx.TimeFmt)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, constx.TimeFmt)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(constx.TimeFmt)
}

func (t *Date) UnmarshalJSON(data []byte) error {
	tt, err := time.ParseInLocation(`"`+constx.TimeFmt+`"`, string(data), time.Local)
	if err != nil {
		return err
	}
	*t = Date(tt)
	return nil
}

func (t Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(constx.DateFmt)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, constx.DateFmt)
	b = append(b, '"')
	return b, nil
}

func (t Date) String() string {
	return time.Time(t).Format(constx.DateFmt)
}
