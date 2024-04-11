package anyx

import (
	"reflect"
)

// if取值
func If[T any](x bool, t, f T) T {
	if x {
		return t
	} else {
		return f
	}
}

// 判空
func IfZeroElse[T any](x, t, f T) T {
	if reflect.ValueOf(&x).Elem().IsZero() {
		return t
	} else {
		return f
	}
}

// 判空取默认值
func IfZero[T any](x, t T) T {
	if reflect.ValueOf(&x).Elem().IsZero() {
		return t
	}
	return x
}

func Default[T any](x []T, d T) T {
	if len(x) == 0 {
		return d
	} else {
		return x[0]
	}
}
