package floatx

import (
	"strconv"
)

func NewFloat64(t ...float64) *Float {
	var x = &Float{notnull: true}
	if len(t) > 0 {
		x.value = t[0]
	} else {
		x.value = 0
	}
	return x
}

type Float struct {
	value   float64
	notnull bool
}

func (x *Float) UnmarshalJSON(bytes []byte) error {
	if str := string(bytes); str != "" && str != "null" {
		if value, err := strconv.ParseFloat(str, 64); err == nil {
			x.value = value
			x.notnull = true
			return nil
		}
	}
	x.notnull = false
	return nil
}

func (x *Float) MarshalJSON() ([]byte, error) {
	if x != nil && x.notnull {
		return []byte(strconv.FormatFloat(x.value, 'f', -1, 64)), nil
	} else {
		return []byte("null"), nil
	}
}

func (x *Float) Value(def ...float64) float64 {
	return x.Float64(def...)
}

func (x *Float) NotNull() bool {
	return x.notnull
}

func (x *Float) String(def ...string) string {
	if x != nil && x.notnull {
		return strconv.FormatFloat(x.value, 'f', -1, 64)
	} else if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (x *Float) Int(def ...int) int {
	if x != nil && x.notnull {
		return int(x.value)
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x *Float) Int64(def ...int64) int64 {
	if x != nil && x.notnull {
		return int64(int(x.value))
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x *Float) Float64(def ...float64) float64 {
	if x != nil && x.notnull {
		return x.value
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x *Float) Bool(def ...bool) bool {
	if x != nil && x.notnull {
		return x.value == 1
	} else if len(def) > 0 {
		return def[0]
	}
	return false
}
