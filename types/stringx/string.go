package stringx

import (
	"strconv"
)

func NewString(v ...string) *String {
	var x = &String{notnull: true}
	if len(v) > 0 {
		x.value = v[0]
	} else {
		x.value = ""
	}
	return x
}

type String struct {
	value   string
	notnull bool
}

func (x *String) UnmarshalJSON(bytes []byte) error {
	if l := len(bytes); l >= 0 && string(bytes) != "null" {
		if l > 1 && bytes[0] == 34 && bytes[l-1] == 34 {
			bytes = bytes[1 : l-1] // 带引号则去掉引号
		}
		x.notnull = true
		x.value = string(bytes)
		return nil
	}
	x.notnull = false
	return nil
}

func (x *String) MarshalJSON() ([]byte, error) {
	if x.Valid() {
		var bytes []byte
		bytes = append(bytes, 34)
		bytes = append(bytes, []byte(x.value)...)
		bytes = append(bytes, 34)
		return bytes, nil
	}
	return []byte("null"), nil
}

func (x *String) Value(def ...string) string {
	return x.String(def...)
}

func (x *String) Valid() bool {
	return x != nil && x.notnull
}

func (x *String) String(def ...string) string {
	if x.Valid() {
		return x.value
	} else if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (x *String) Int(def ...int) int {
	if x.Valid() {
		if value, err := strconv.Atoi(x.value); err == nil {
			return value
		}
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x *String) Int64(def ...int64) int64 {
	if x.Valid() {
		if value, err := strconv.ParseInt(x.value, 10, 64); err == nil {
			return value
		}
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x *String) Float64(def ...float64) float64 {
	if x.Valid() {
		if value, err := strconv.ParseFloat(x.value, 64); err == nil {
			return value
		}
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x *String) Bool(def ...bool) bool {
	if x.Valid() {
		return ParseBool(x.value)
	} else if len(def) > 0 {
		return def[0]
	}
	return false
}
