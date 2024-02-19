package anyx

import (
	"reflect"
)

// if执行方法
func IfElseFunc(exp bool, trueFunc, falseFunc func()) {
	if exp {
		trueFunc()
	} else {
		falseFunc()
	}
}

// if取值
func IfElseValue[T any](exp bool, trueValue, falseValue T) T {
	if exp {
		return trueValue
	} else {
		return falseValue
	}
}

// 判空
func IfElseZero[T any](target, defaultValue, elseValue T) T {
	if reflect.ValueOf(&target).Elem().IsZero() {
		return defaultValue
	} else {
		return elseValue
	}
}

// 判空取默认值
func IfZero[T any](target, defaultValue T) T {
	if reflect.ValueOf(&target).Elem().IsZero() {
		return defaultValue
	}
	return target
}
