package structx

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/magiconair/properties"
	"gopkg.in/yaml.v3"

	"github.com/go-xuan/quanx/utils/filex"
)

// 写入结构体到json文件
func WriteStructToJson(filePath string, obj interface{}) (err error) {
	var bytes []byte
	bytes, err = json.MarshalIndent(obj, "", "	")
	if err != nil {
		return
	}
	err = os.WriteFile(filePath, bytes, 0777)
	if err != nil {
		return
	}
	return
}

// 读取json文件到结构体
func ReadJsonToStruct(filePath string, obj interface{}) (err error) {
	var bytes []byte
	bytes, err = os.ReadFile(filePath)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, obj)
	if err != nil {
		return
	}
	return
}

func SetDefaultValue(config interface{}) error {
	valueRef := reflect.ValueOf(config)
	if valueRef.Type().Kind() != reflect.Ptr {
		// 修改值必须是指针类型否则不可行
		return errors.New("config must be pointer type")
	}
	for i := 0; i < valueRef.Elem().NumField(); i++ {
		field := valueRef.Elem().Field(i)
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
	return nil
}

// 读取配置文件到指针
func ReadConfigToPointer(filePath string, config interface{}) (err error) {
	var bytes []byte
	bytes, err = os.ReadFile(filePath)
	if err != nil {
		return
	}
	switch filex.Suffix(filePath) {
	case filex.Json:
		err = ReadJsonToPointer(bytes, config)
	case filex.Yaml, filex.Yml:
		err = ReadYamlToPointer(bytes, config)
	case filex.Toml:
		err = ReadTomlToPointer(bytes, config)
	case filex.Properties:
		err = ReadPropertiesToPointer(bytes, config)
	}
	return
}

// 读取json文件到指针
func ReadJsonToPointer(bytes []byte, config interface{}) (err error) {
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return
	}
	return
}

// 读取yaml文件到指针
func ReadYamlToPointer(bytes []byte, config interface{}) (err error) {
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return
	}
	return
}

// 读取toml文件到指针
func ReadTomlToPointer(bytes []byte, config interface{}) (err error) {
	err = toml.Unmarshal(bytes, config)
	if err != nil {
		return
	}
	return
}

// 读取toml文件到指针
func ReadPropertiesToPointer(bytes []byte, config interface{}) (err error) {
	valueRef := reflect.ValueOf(config)
	if valueRef.Type().Kind() != reflect.Ptr {
		// 修改值必须是指针类型否则不可行
		return errors.New("config must be pointer type")
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
		default:
			if tag == "" {
				continue
			}
			fmt.Printf("未匹配到type %s", field.Kind())
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
		default:
			if tag == "" {
				continue
			}
			fmt.Printf("未匹配到type %s", field.Kind())
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
		default:
			if tag == "" {
				continue
			}
			fmt.Printf("未匹配到type %s", field.Kind())
		}
	}
}
