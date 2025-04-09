package intx

import (
	"strconv"
)

func NewInt64(v ...int64) Int64 {
	var x = Int64{notnull: true}
	if len(v) > 0 {
		x.value = v[0]
	} else {
		x.value = 0
	}
	return x
}

type Int64 struct {
	value   int64
	notnull bool
}

func (x *Int64) UnmarshalJSON(bytes []byte) error {
	if value, err := strconv.ParseInt(string(bytes), 10, 64); err != nil {
		return err
	} else {
		x.value = value
		x.notnull = true
	}
	return nil
}

func (x Int64) MarshalJSON() ([]byte, error) {
	if x.notnull {

		return []byte(strconv.FormatInt(x.value, 10)), nil
	} else {
		return []byte("null"), nil
	}
}

func (x Int64) Value(def ...int64) int64 {
	return x.Int64(def...)
}

func (x Int64) NotNull() bool {
	return x.notnull
}

func (x Int64) String(def ...string) string {
	if x.notnull {
		return strconv.FormatInt(x.value, 10)
	} else if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (x Int64) Int(def ...int) int {
	if x.notnull {
		return int(x.value)
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x Int64) Int64(def ...int64) int64 {
	if x.notnull {
		return x.value
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x Int64) Float64(def ...float64) float64 {
	if x.notnull {
		return float64(x.value)
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x Int64) Bool(def ...bool) bool {
	if x.notnull {
		return x.value == 1
	} else if len(def) > 0 {
		return def[0]
	}
	return false
}
