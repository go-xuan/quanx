package propertiesx

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/magiconair/properties"
)

// 读取properties文件到指针
func Unmarshal(bytes []byte, config interface{}) (err error) {
	valueRef := reflect.ValueOf(config)
	if valueRef.Type().Kind() != reflect.Ptr {
		// 对象必须是指针类型
		return errors.New("the config must be pointer type")
	}
	var pp *properties.Properties
	pp, err = properties.Load(bytes, properties.UTF8)
	if err != nil {
		return
	}
	for i := 0; i < valueRef.Elem().NumField(); i++ {
		field := valueRef.Elem().Field(i)
		tag := valueRef.Elem().Type().Field(i).Tag.Get("properties")
		switch field.Kind() {
		case reflect.String:
			field.SetString(pp.GetString(tag, ""))
		case reflect.Bool:
			field.SetBool(pp.GetBool(tag, false))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(int64(pp.GetInt(tag, 0)))
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.SetInt(int64(pp.GetInt(tag, 0)))
		case reflect.Float32, reflect.Float64:
			field.SetFloat(pp.GetFloat64(tag, 0))
		case reflect.Struct:
			setStructValue(pp, field)
		case reflect.Pointer:
			setPointerValue(pp, field)
		default:
			if tag == "" {
				continue
			}
			fmt.Println("the type not matched : ", field.Kind())
		}
	}
	return
}

func setStructValue(pp *properties.Properties, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := v.Type().Field(i).Tag.Get("properties")
		switch field.Kind() {
		case reflect.String:
			field.SetString(pp.GetString(tag, ""))
		case reflect.Bool:
			field.SetBool(pp.GetBool(tag, false))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(int64(pp.GetInt(tag, 0)))
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.SetInt(int64(pp.GetInt(tag, 0)))
		case reflect.Float32, reflect.Float64:
			field.SetFloat(pp.GetFloat64(tag, 0))
		case reflect.Struct:
			setStructValue(pp, field)
		case reflect.Pointer:
			setPointerValue(pp, field)
		default:
			if tag == "" {
				continue
			}
			fmt.Println("the type not matched : ", field.Kind())
		}
	}
}

func setPointerValue(pp *properties.Properties, v reflect.Value) {
	for i := 0; i < v.Elem().NumField(); i++ {
		field := v.Elem().Field(i)
		tag := v.Elem().Type().Field(i).Tag.Get("properties")
		switch field.Kind() {
		case reflect.String:
			field.SetString(pp.GetString(tag, ""))
		case reflect.Bool:
			field.SetBool(pp.GetBool(tag, false))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(int64(pp.GetInt(tag, 0)))
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.SetInt(int64(pp.GetInt(tag, 0)))
		case reflect.Float32, reflect.Float64:
			field.SetFloat(pp.GetFloat64(tag, 0))
		case reflect.Struct:
			setStructValue(pp, field)
		case reflect.Pointer:
			setPointerValue(pp, field)
		default:
			if tag == "" {
				continue
			}
			fmt.Println("the type not matched : ", field.Kind())
		}
	}
}
