package anyx

import (
	"fmt"
	"strconv"
	"strings"
)

type Value interface {
	String() string
	Int() int
	Int64() int64
	Float64() float64
	Bool() bool
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
	return nil
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

func (v *stringValue) String() string {
	return v.val
}

func (v *stringValue) Int() int {
	if value, err := strconv.Atoi(v.val); err != nil {
		return 0
	} else {
		return value
	}
}

func (v *stringValue) Int64() int64 {
	if value, err := strconv.ParseInt(v.val, 10, 64); err != nil {
		return 0
	} else {
		return value
	}
}

func (v *stringValue) Float64() float64 {
	if value, err := strconv.ParseFloat(v.val, 64); err != nil {
		return 0
	} else {
		return value
	}
}

func (v *stringValue) Bool() bool {
	switch strings.ToLower(v.val) {
	case "true", "是", "yes":
		return true
	default:
		return false
	}
}

type intValue struct {
	val int
}

func (v *intValue) String() string {
	return strconv.Itoa(v.val)
}

func (v *intValue) Int() int {
	return v.val
}

func (v *intValue) Int64() int64 {
	return int64(v.val)
}

func (v *intValue) Float64() float64 {
	return float64(v.val)
}

func (v *intValue) Bool() bool {
	if v.val == 1 {
		return true
	} else {
		return false
	}
}

type int64Value struct {
	val int64
}

func (v *int64Value) String() string {
	return strconv.FormatInt(v.val, 10)
}

func (v *int64Value) Int() int {
	return int(v.val)
}

func (v *int64Value) Int64() int64 {
	return v.val
}

func (v *int64Value) Float64() float64 {
	return float64(v.val)
}

func (v *int64Value) Bool() bool {
	if v.val == 1 {
		return true
	} else {
		return false
	}
}

type float64Value struct {
	val float64
}

func (v *float64Value) String() string {
	return strconv.FormatFloat(v.val, 'f', -1, 64)
}

func (v *float64Value) Int() int {
	return int(v.val)
}

func (v *float64Value) Int64() int64 {
	return int64(v.val)
}

func (v *float64Value) Float64() float64 {
	return v.val
}

func (v *float64Value) Bool() bool {
	if v.val == 1 {
		return true
	} else {
		return false
	}
}

type boolValue struct {
	val bool
}

func (v *boolValue) String() string {
	if v.val {
		return "true"
	}
	return "false"
}

func (v *boolValue) Int() int {
	if v.val {
		return 1
	}
	return 0
}

func (v *boolValue) Int64() int64 {
	if v.val {
		return 1
	}
	return 0
}

func (v *boolValue) Float64() float64 {
	if v.val {
		return 1
	}
	return 0
}

func (v *boolValue) Bool() bool {
	return v.val
}
