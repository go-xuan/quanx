package osx

import (
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/boolx"
	"github.com/go-xuan/quanx/types/floatx"
	"github.com/go-xuan/quanx/types/intx"
	"github.com/go-xuan/quanx/types/stringx"
)

const (
	WindowsOS = "windows"
	LinuxOS   = "linux"
	MacOS     = "darwin"
)

func GetHome() string {
	switch runtime.GOOS {
	case WindowsOS:
		return os.Getenv("USERPROFILE")
	case MacOS:
		return os.Getenv("HOME")
	default:
		return ""
	}
}

func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return os.Getenv("HOSTNAME")
	}
	return hostname
}

// SetValueFromEnv 从环境变量中读取值
func SetValueFromEnv(v any) error {
	if err := anyx.MustStructPointer(v); err != nil {
		return errorx.New("the kind must be struct pointer")
	}
	var elem = reflect.ValueOf(v).Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		tag := elem.Type().Field(i).Tag.Get("env")
		key := strings.ToUpper(tag)
		value := stringx.NewString(os.Getenv(key))
		switch field.Kind() {
		case reflect.String:
			field.SetString(value.Value())
		case reflect.Bool:
			field.SetBool(value.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.SetInt(value.Int64())
		case reflect.Float32, reflect.Float64:
			field.SetFloat(value.Float64())
		default:
			continue
		}
	}
	return nil
}

// SetEnvFromValue 将值更新到环境变量
func SetEnvFromValue(v any) error {
	if err := anyx.MustStructPointer(v); err != nil {
		return errorx.New("the kind must be struct pointer")
	}
	var elem = reflect.ValueOf(v).Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		tag := elem.Type().Field(i).Tag.Get("env")
		key := strings.ToUpper(tag)
		var value string
		switch field.Kind() {
		case reflect.String:
			value = field.String()
		case reflect.Bool:
			value = boolx.NewBool(field.Bool()).String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = intx.NewInt64(field.Int()).String()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = intx.NewInt64(int64(field.Uint())).String()
		case reflect.Float32, reflect.Float64:
			value = floatx.NewFloat64(field.Float()).String()
		default:
			continue // 忽略不支持的类型
		}
		if err := os.Setenv(key, value); err != nil {
			return errorx.Wrap(err, "set environment variable error")
		}
	}
	return nil
}
