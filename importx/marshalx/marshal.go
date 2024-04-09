package marshalx

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/go-xuan/quanx/utilx/filex"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v3"
)

type CaseType uint

type Case struct {
	marshal   func(any) ([]byte, error)
	unmarshal func([]byte, any) error
}

func Marshal(t string, v any) ([]byte, error) {
	return GetCase(t).marshal(v)
}

func Unmarshal(t string, b []byte, v any) error {
	return GetCase(t).unmarshal(b, v)
}

func GetCase(t string) *Case {
	switch t {
	case filex.Json:
		return &Case{marshal: json.Marshal, unmarshal: json.Unmarshal}
	case filex.Yaml, filex.Yml:
		return &Case{marshal: yaml.Marshal, unmarshal: yaml.Unmarshal}
	case filex.Toml:
		return &Case{marshal: TomlMarshal, unmarshal: toml.Unmarshal}
	case filex.Properties:
		return &Case{marshal: PropertiesMarshal, unmarshal: PropertiesUnmarshal}
	case "msgpack":
		return &Case{marshal: msgpack.Marshal, unmarshal: msgpack.Unmarshal}
	default:
		return &Case{marshal: DefaultMarshal, unmarshal: DefaultUnmarshal}
	}
}

func DefaultMarshal(v any) (out []byte, err error) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.String:
		out = []byte(v.(string))
	case reflect.Bool:
		out = []byte(strconv.FormatBool(v.(bool)))
	case reflect.Int:
		out = []byte(strconv.FormatInt(int64(v.(int)), 10))
	case reflect.Int64:
		out = []byte(strconv.FormatInt(v.(int64), 10))
	case reflect.Float32:
		out = []byte(strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32))
	case reflect.Float64:
		out = []byte(strconv.FormatFloat(v.(float64), 'f', -1, 64))
	default:
		err = errors.New("unsupported type")
	}
	return
}

func DefaultUnmarshal(data []byte, v any) (err error) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.String:
		v = string(data)
	case reflect.Bool:
		v = string(data) == "true"
	case reflect.Int:
		v, err = strconv.Atoi(string(data))
	case reflect.Int64:
		v, err = strconv.ParseInt(string(data), 10, 64)
	case reflect.Float32:
		v, err = strconv.ParseFloat(string(data), 32)
	case reflect.Float64:
		v, err = strconv.ParseFloat(string(data), 64)
	default:
		err = errors.New("unsupported type")
	}
	return
}

// 读取配置文件到指针
func LoadFromFile(path string, v any) (err error) {
	if !filex.Exists(path) {
		return errors.New("the file not exist : " + path)
	}
	var b []byte
	if b, err = os.ReadFile(path); err != nil {
		return
	}
	var suffix = filex.Suffix(path)
	if err = GetCase(suffix).unmarshal(b, v); err != nil {
		return
	}
	return
}
