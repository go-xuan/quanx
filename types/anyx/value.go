package anyx

import (
	"fmt"
	"strconv"
	"strings"
)

type Value interface {
	String(def ...string) string
	Int(def ...int) int
	Int64(def ...int64) int64
	Float64(def ...float64) float64
	Bool(def ...bool) bool
}

func ValueOf(v any) Value {
	switch value := v.(type) {
	case int:
		return IntValue(value)
	case int64:
		return Int64Value(value)
	case float64:
		return Float64Value(value)
	case bool:
		return BoolValue(value)
	case string:
		return StringValue(value)
	default:
		panic(fmt.Sprintf("unsupported value type: %T", value))
	}
	return ZeroValue()
}

func ZeroValue() Value {
	return StringValue("")
}

func StringValue(v string) Value {
	return &stringValue{v}
}

func Int64Value(v int64) Value {
	return &int64Value{v}
}

func IntValue(v int) Value {
	return &intValue{v}
}

func Float64Value(v float64) Value {
	return &float64Value{v}
}

func BoolValue(v bool) Value {
	return &boolValue{v}
}

type stringValue struct {
	val string
}

func (v *stringValue) String(def ...string) string {
	if v.val == "" && len(def) > 0 {
		return def[0]
	}
	return v.val
}

func (v *stringValue) Int(def ...int) int {
	if value, err := strconv.Atoi(v.val); err != nil && len(def) > 0 {
		return def[0]
	} else {
		return value
	}
}

func (v *stringValue) Int64(def ...int64) int64 {
	if value, err := strconv.ParseInt(v.val, 10, 64); err != nil && len(def) > 0 {
		return def[0]
	} else {
		return value
	}
}

func (v *stringValue) Float64(def ...float64) float64 {
	if value, err := strconv.ParseFloat(v.val, 64); err != nil && len(def) > 0 {
		return def[0]
	} else {
		return value
	}
}

func (v *stringValue) Bool(def ...bool) bool {
	switch strings.ToLower(v.val) {
	case "true", "是", "yes":
		return true
	case "false", "否", "no":
		return false
	default:
		if len(def) > 0 {
			return def[0]
		}
		return false
	}
}

type intValue struct {
	val int
}

func (v *intValue) String(def ...string) string {
	if value := strconv.Itoa(v.val); value == "" && len(def) > 0 {
		return def[0]
	} else {
		return value
	}
}

func (v *intValue) Int(def ...int) int {
	if v.val == 0 && len(def) > 0 {
		return def[0]
	}
	return v.val
}

func (v *intValue) Int64(def ...int64) int64 {
	if value := int64(v.val); value == 0 && len(def) > 0 {
		return def[0]
	} else {
		return value
	}
}

func (v *intValue) Float64(def ...float64) float64 {
	if value := float64(v.val); value == 0 && len(def) > 0 {
		return def[0]
	} else {
		return value
	}
}

func (v *intValue) Bool(def ...bool) bool {
	if v.val == 1 {
		return true
	} else if len(def) > 0 {
		return def[0]
	} else {
		return false
	}
}

type int64Value struct {
	val int64
}

func (v *int64Value) String(def ...string) string {
	if value := strconv.FormatInt(v.val, 10); value == "" && len(def) > 0 {
		return def[0]
	} else {
		return value
	}
}

func (v *int64Value) Int(def ...int) int {
	if value := int(v.val); value == 0 && len(def) > 0 {
		return def[0]
	} else {
		return value
	}
}

func (v *int64Value) Int64(def ...int64) int64 {
	if v.val == 0 && len(def) > 0 {
		return def[0]
	}
	return v.val
}

func (v *int64Value) Float64(def ...float64) float64 {
	if value := float64(v.val); value == 0 && len(def) > 0 {
		return def[0]
	} else {
		return value
	}
}

func (v *int64Value) Bool(def ...bool) bool {
	if v.val == 1 {
		return true
	} else if len(def) > 0 {
		return def[0]
	} else {
		return false
	}
}

type float64Value struct {
	val float64
}

func (v *float64Value) String(def ...string) string {
	if value := strconv.FormatFloat(v.val, 'f', -1, 64); value == "" && len(def) > 0 {
		return def[0]
	} else {
		return value
	}
}

func (v *float64Value) Int(def ...int) int {
	if value := int(v.val); value == 0 && len(def) > 0 {
		return def[0]
	} else {
		return value
	}
}

func (v *float64Value) Int64(def ...int64) int64 {
	if value := int64(v.val); value == 0 && len(def) > 0 {
		return def[0]
	} else {
		return value
	}
}

func (v *float64Value) Float64(def ...float64) float64 {
	if v.val == 0 && len(def) > 0 {
		return def[0]
	} else {
		return v.val
	}
}

func (v *float64Value) Bool(def ...bool) bool {
	if v.val == 1 {
		return true
	} else if len(def) > 0 {
		return def[0]
	} else {
		return false
	}
}

type boolValue struct {
	val bool
}

func (v *boolValue) String(def ...string) string {
	if v.val {
		return "true"
	} else if len(def) > 0 {
		return def[0]
	} else {
		return "false"
	}
}

func (v *boolValue) Int(def ...int) int {
	if v.val {
		return 1
	} else if len(def) > 0 {
		return def[0]
	} else {
		return 0
	}
}

func (v *boolValue) Int64(def ...int64) int64 {
	if v.val {
		return 1
	} else if len(def) > 0 {
		return def[0]
	} else {
		return 0
	}
}

func (v *boolValue) Float64(def ...float64) float64 {
	if v.val {
		return 1
	} else if len(def) > 0 {
		return def[0]
	} else {
		return 0
	}
}

func (v *boolValue) Bool(def ...bool) bool {
	if v.val {
		return true
	} else if len(def) > 0 {
		return def[0]
	} else {
		return false
	}
}
