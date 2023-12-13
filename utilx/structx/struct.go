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

	"github.com/go-xuan/quanx/utilx/filex"
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

// 设置默认值
func SetDefaultValue(config interface{}) error {
	valueRef := reflect.ValueOf(config)
	if valueRef.Type().Kind() != reflect.Ptr {
		// 对象必须是指针类型
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

// 读取配置文件到指针
func ReadFileToPointer(config interface{}, filePath string) (err error) {
	var bytes []byte
	bytes, err = os.ReadFile(filePath)
	if err != nil {
		return
	}
	err = ParseBytesToPointer(config, bytes, filex.Suffix(filePath))
	if err != nil {
		return
	}
	return
}

// 解析bytes到指针
func ParseBytesToPointer(config interface{}, bytes []byte, filePath string) (err error) {
	var fileType = filex.Suffix(filePath)
	switch fileType {
	case filex.Json:
		err = ReadJsonToPointer(config, bytes)
	case filex.Yaml, filex.Yml:
		err = ReadYamlToPointer(config, bytes)
	case filex.Toml:
		err = ReadTomlToPointer(config, bytes)
	case filex.Properties:
		err = ReadPropertiesToPointer(config, bytes)
	default:
		err = errors.New("当前文件类型暂不支持:" + fileType)
	}
	if err != nil {
		return
	}
	return
}

// 读取bytes到指针
func ReadJsonToPointer(config interface{}, bytes []byte) (err error) {
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return
	}
	return
}

// 读取yaml文件到指针
func ReadYamlToPointer(config interface{}, bytes []byte) (err error) {
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return
	}
	return
}

// 读取toml文件到指针
func ReadTomlToPointer(config interface{}, bytes []byte) (err error) {
	err = toml.Unmarshal(bytes, config)
	if err != nil {
		return
	}
	return
}

// 读取properties文件到指针
func ReadPropertiesToPointer(config interface{}, bytes []byte) (err error) {
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
