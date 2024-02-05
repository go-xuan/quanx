package anyx

import (
	"errors"
	"reflect"
	"strconv"
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

// 设置默认值
func SetDefaultValue[T any](config T) error {
	valueRef := reflect.ValueOf(config)
	if valueRef.Type().Kind() != reflect.Ptr {
		return errors.New("the config must be pointer type")
	}
	for i := 0; i < valueRef.Elem().NumField(); i++ {
		field := valueRef.Elem().Field(i)
		if field.IsZero() {
			value := valueRef.Elem().Type().Field(i).Tag.Get("default")
			switch field.Kind() {
			case reflect.Bool:
				field.SetBool(value == "true")
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				intVal, _ := strconv.ParseInt(value, 10, 64)
				field.SetInt(intVal)
			case reflect.Float32, reflect.Float64:
				floatValue, _ := strconv.ParseFloat(value, 64)
				field.SetFloat(floatValue)
			default:
				field.SetString(value)
			}
		}
	}
	return nil
}
