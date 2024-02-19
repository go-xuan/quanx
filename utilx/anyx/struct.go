package anyx

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// 设置默认值
func SetDefaultValue(obj interface{}) error {
	valueRef := reflect.ValueOf(obj)
	if valueRef.Type().Kind() != reflect.Ptr {
		return errors.New("the obj must be pointer type")
	}
	for i := 0; i < valueRef.Elem().NumField(); i++ {
		field := valueRef.Elem().Field(i)
		if field.IsZero() {
			value := valueRef.Elem().Type().Field(i).Tag.Get("default")
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
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
				return fmt.Errorf("unsupported type %T", value)
			}
		}
	}
	return nil
}

func MapToStruct(m map[string]string, obj interface{}) error {
	v := reflect.ValueOf(obj).Elem() // 获取指向结构体的值类型
	for key, value := range m {
		field := v.FieldByName(key) // 根据字段名称查找对应的字段
		if field.IsValid() && field.CanSet() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
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
			return fmt.Errorf("no such field %q in the structure", key)
		}
	}
	return nil
}
