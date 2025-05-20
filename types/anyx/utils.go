package anyx

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/types/boolx"
	"github.com/go-xuan/quanx/types/stringx"
)

//  获取结构体的值
func getValueOf(v any) reflect.Value {
	var valueOf = reflect.ValueOf(v)
	if valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	return valueOf
}

// GetTagNames 获取结构体中指定tag名
func GetTagNames(v any, tag string) []string {
	valueOf := getValueOf(v)
	typeOf := valueOf.Type()

	var names []string
	for i := 0; i < typeOf.NumField(); i++ {
		if name, ok := typeOf.Field(i).Tag.Lookup(tag); ok {
			names = append(names, name)
		}
	}
	return names
}

// GetTagFieldValues 获取结构体中指定tag的字段值
func GetTagFieldValues(v any, tag string) []string {
	valueOf := getValueOf(v)
	typeOf := valueOf.Type()

	var values []string
	for i := 0; i < typeOf.NumField(); i++ {
		if _, ok := typeOf.Field(i).Tag.Lookup(tag); ok {
			values = append(values, fmt.Sprintf("%v", valueOf.Field(i).Interface()))
		}
	}
	return values
}

// MustStructPointer 检查是否为结构体指针
func MustStructPointer(v any) error {
	var typeOf = reflect.TypeOf(v)
	if typeOf.Kind() != reflect.Pointer {
		return errorx.New("the kind must be pointer")
	} else if typeOf.Elem().Kind() == reflect.Slice {
		return errorx.New("the kind cannot be slice")
	}
	return nil
}

// SetDefaultValue 设置默认值
func SetDefaultValue(v interface{}, tag ...string) error {
	if err := MustStructPointer(v); err != nil {
		return errorx.New("the kind must be struct pointer")
	}
	var elem = reflect.ValueOf(v).Elem()
	key := stringx.Default("default", tag...)
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if field.IsZero() {
			if value := elem.Type().Field(i).Tag.Get(key); value != "" {
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

// MapToStruct 将map转换为结构体
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

// MergeStructs 合并结构体
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

// SetZeroValue 设置结构体零值
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
