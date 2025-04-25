package boolx

import (
	"strconv"
	"strings"
)

func NewBool(v ...bool) Bool {
	var x = Bool{notnull: true}
	if len(v) > 0 && v[0] {
		x.value = true
	}
	return x
}

type Bool struct {
	value   bool
	notnull bool
}

func (x *Bool) UnmarshalJSON(bytes []byte) error {
	if value := string(bytes); value != "" {
		x.notnull = true
		x.value = ValueOf(value)
	} else {
		x.notnull = false
	}
	return nil
}

func (x Bool) MarshalJSON() ([]byte, error) {
	if x.notnull && x.value {
		return []byte("true"), nil
	} else {
		return []byte("false"), nil
	}
}

func (x Bool) Value(def ...bool) bool {
	return x.Bool(def...)
}

func (x Bool) NotNull() bool {
	return x.notnull
}

func (x Bool) String(def ...string) string {
	if x.notnull {
		return strconv.FormatBool(x.value)
	} else if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (x Bool) Int(def ...int) int {
	if x.notnull && x.value {
		return 1
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x Bool) Int64(def ...int64) int64 {
	if x.notnull && x.value {
		return 1
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x Bool) Float64(def ...float64) float64 {
	if x.notnull && x.value {
		return 1
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x Bool) Bool(def ...bool) bool {
	if x.notnull {
		return x.value
	} else if len(def) > 0 {
		return def[0]
	}
	return false
}

func ValueOf(s string, def ...bool) bool {
	switch strings.ToLower(s) {
	case "1", "t", "true", "yes", "是":
		return true
	default:
		if len(def) > 0 {
			return def[0]
		} else {
			return false
		}
	}
}
