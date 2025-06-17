package anyx

import (
	"github.com/go-xuan/quanx/types/boolx"
	"github.com/go-xuan/quanx/types/floatx"
	"github.com/go-xuan/quanx/types/intx"
	"github.com/go-xuan/quanx/types/stringx"
)

type Value interface {
	Valid() bool
	String(def ...string) string
	Int(def ...int) int
	Int64(def ...int64) int64
	Float64(def ...float64) float64
	Bool(def ...bool) bool
}

type Zero struct{}

func (z *Zero) Valid() bool {
	return false
}

func (z *Zero) String(...string) string {
	return ""
}

func (z *Zero) Int(...int) int {
	return 0
}

func (z *Zero) Int64(...int64) int64 {
	return 0
}

func (z *Zero) Float64(...float64) float64 {
	return 0
}

func (z *Zero) Bool(...bool) bool {
	return false
}

func New(v any) Value {
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
		return ZeroValue()
	}
}

func ZeroValue() Value {
	return &Zero{}
}

func StringValue(v string) Value {
	return stringx.NewString(v)
}

func Int64Value(v int64) Value {
	return intx.NewInt64(v)
}

func IntValue(v int) Value {
	return intx.NewInt(v)
}

func Float64Value(v float64) Value {
	return floatx.NewFloat64(v)
}

func BoolValue(v bool) Value {
	return boolx.NewBool(v)
}
