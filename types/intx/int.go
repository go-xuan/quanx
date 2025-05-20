package intx

import "strconv"

func NewInt(v ...int) Int {
	var x = Int{notnull: true}
	if len(v) > 0 {
		x.value = v[0]
	} else {
		x.value = 0
	}
	return x
}

type Int struct {
	value   int
	notnull bool
}

func (x *Int) UnmarshalJSON(bytes []byte) error {
	if str := string(bytes); str != "" && str != "null" {
		if value, err := strconv.Atoi(str); err == nil {
			x.value = value
			x.notnull = true
			return nil
		}
	}
	x.notnull = false
	return nil
}

func (x *Int) MarshalJSON() ([]byte, error) {
	if x != nil && x.notnull {
		return []byte(strconv.Itoa(x.value)), nil
	}
	return []byte("null"), nil
}

func (x *Int) Value(def ...int) int {
	return x.Int(def...)
}

func (x *Int) NotNull() bool {
	return x.notnull
}

func (x *Int) String(def ...string) string {
	if x != nil && x.notnull {
		return strconv.Itoa(x.value)
	} else if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (x *Int) Int(def ...int) int {
	if x != nil && x.notnull {
		return x.value
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x *Int) Int64(def ...int64) int64 {
	if x != nil && x.notnull {
		return int64(x.value)
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x *Int) Float64(def ...float64) float64 {
	if x != nil && x.notnull {
		return float64(x.value)
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (x *Int) Bool(def ...bool) bool {
	if x != nil && x.notnull {
		return x.value == 1
	} else if len(def) > 0 {
		return def[0]
	}
	return false
}
