package anyx

import "reflect"

// if取值
func IfValue[T any](exp bool, ifV, elseV T) T {
	if exp {
		return ifV
	} else {
		return elseV
	}
}

// if执行方法
func IfFunc(exp bool, ifF, elseF func()) {
	if exp {
		ifF()
	} else {
		elseF()
	}
}

// 判空取默认值
func IfZero[T any](target, def T) T {
	if reflect.ValueOf(&target).Elem().IsZero() {
		return def
	}
	return target
}
