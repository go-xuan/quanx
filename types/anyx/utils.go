package anyx

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	
	"github.com/go-xuan/quanx/types/boolx"
)

// SetDefaultValue 设置默认值
func SetDefaultValue(v interface{}, tag ...string) error {
	valueRef, key := reflect.ValueOf(v), "default"
	if len(tag) > 0 && tag[0] != "" {
		key = tag[0]
	}
	if valueRef.Type().Kind() != reflect.Ptr {
		return errors.New("param must be pointer type")
	}
	for i := 0; i < valueRef.Elem().NumField(); i++ {
		field := valueRef.Elem().Field(i)
		if field.IsZero() {
			if value := valueRef.Elem().Type().Field(i).Tag.Get(key); value != "" {
				switch field.Kind() {
				case reflect.String:
					field.SetString(value)
				case reflect.Bool:
					field.SetBool(boolx.ValueOf(value))
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					intVal, _ := strconv.ParseInt(value, 10, 64)
					field.SetInt(intVal)
				case reflect.Float32, reflect.Float64:
					floatValue, _ := strconv.ParseFloat(value, 64)
					field.SetFloat(floatValue)
				default:
					return fmt.Errorf("unsupported type %T", value)
				}
			}
		}
	}
	return nil
}

func MapToStruct(m map[string]string, v interface{}) error {
	elem := reflect.ValueOf(v).Elem() // 获取指向结构体的值类型
	for key, value := range m {
		field := elem.FieldByName(key) // 根据字段名称查找对应的字段
		if field.IsValid() && field.CanSet() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Bool:
				field.SetBool(boolx.ValueOf(value))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				intVal, _ := strconv.ParseInt(value, 10, 64)
				field.SetInt(intVal)
			case reflect.Float32, reflect.Float64:
				floatValue, _ := strconv.ParseFloat(value, 64)
				field.SetFloat(floatValue)
			default:
				return fmt.Errorf("unsupported type %T", value)
			}
		} else if !field.IsValid() {
			// 如果没有找到对应的字段则返回错误
			return fmt.Errorf("no such field %s in the structure", key)
		}
	}
	return nil
}

func MergeStructs(a, b interface{}) {
	va, vb := reflect.ValueOf(a).Elem(), reflect.ValueOf(b).Elem()
	for i := 0; i < va.NumField(); i++ {
		fieldA := va.Type().Field(i)
		fieldB := vb.FieldByName(fieldA.Name)
		if fieldB.IsValid() && fieldB.CanSet() {
			va.Field(i).Set(fieldB)
		}
	}
}

func SetZeroValue[T interface{}](a, b T) {
	va, vb := reflect.ValueOf(a).Elem(), reflect.ValueOf(b).Elem()
	for i := 0; i < va.NumField(); i++ {
		if va.Field(i).IsZero() {
			va.Field(i).Set(vb.Field(i))
		}
	}
}

// If if取值
func If[T any](x bool, t, f T) T {
	if x {
		return t
	} else {
		return f
	}
}

// IfZero 判空时取默认值
func IfZero[T any](x, def T) T {
	if reflect.ValueOf(&x).Elem().IsZero() {
		return def
	}
	return x
}

// Default 用于函数中的不定参数取默认值
func Default[T any](def T, x ...T) T {
	if len(x) == 0 {
		return def
	} else {
		return x[0]
	}
}
