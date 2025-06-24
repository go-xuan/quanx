package gormx

import (
	"time"
	
	"github.com/go-xuan/quanx/utils/timex"
)

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	var bytes []byte
	bytes = append(bytes, 34)
	bytes = time.Time(t).AppendFormat(bytes, timex.TimeFmt)
	bytes = append(bytes, 34)
	return bytes, nil
}

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	var bytes []byte
	bytes = append(bytes, 34)
	bytes = time.Time(d).AppendFormat(bytes, timex.DateFmt)
	bytes = append(bytes, 34)
	return bytes, nil
}
