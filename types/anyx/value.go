package anyx

import (
	"github.com/go-xuan/quanx/types/boolx"
	"github.com/go-xuan/quanx/types/floatx"
	"github.com/go-xuan/quanx/types/intx"
	"github.com/go-xuan/quanx/types/stringx"
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
		return ZeroValue()
	}
}

func ZeroValue() Value {
	return IntValue(0)
}

func StringValue(v string) Value {
	x := stringx.NewString(v)
	return &x
}

func Int64Value(v int64) Value {
	x := intx.NewInt64(v)
	return &x
}

func IntValue(v int) Value {
	x := intx.NewInt(v)
	return &x
}

func Float64Value(v float64) Value {
	x := floatx.NewFloat64(v)
	return &x
}

func BoolValue(v bool) Value {
	x := boolx.NewBool(v)
	return &x
}
