package marshalx

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v3"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/filex"
)

// 序列化类型

const (
	Json       = "json"
	Yml        = "yml"
	Yaml       = "yaml"
	Toml       = "toml"
	Properties = "properties"
	Msgpack    = "msgpack"
)

// 序列化方案
type Case struct {
	Name      string
	Marshal   func(any) ([]byte, error)
	Unmarshal func([]byte, any) error
}

// 类型识别
func Identify(s string) string {
	if stringx.ContainsAny(s, ".", "\\", "/") {
		return filex.Suffix(s)
	}
	return s
}

func NewCase(s string) *Case {
	switch Identify(s) {
	case Json:
		return &Case{Name: Json, Marshal: json.Marshal, Unmarshal: json.Unmarshal}
	case Yaml, Yml:
		return &Case{Name: Yaml, Marshal: yaml.Marshal, Unmarshal: yaml.Unmarshal}
	case Toml:
		return &Case{Name: Toml, Marshal: TomlMarshal, Unmarshal: toml.Unmarshal}
	case Properties:
		return &Case{Name: Properties, Marshal: PropertiesMarshal, Unmarshal: PropertiesUnmarshal}
	case Msgpack:
		return &Case{Name: Msgpack, Marshal: msgpack.Marshal, Unmarshal: msgpack.Unmarshal}
	default:
		return DefaultCase()
	}
}

func DefaultCase() *Case {
	return &Case{Name: constx.Default, Marshal: DefaultMarshal, Unmarshal: DefaultUnmarshal}
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
func UnmarshalFromFile(path string, v any) (err error) {
	if !filex.Exists(path) {
		return errors.New("the file not exist : " + path)
	}
	var b []byte
	if b, err = os.ReadFile(path); err != nil {
		return
	}
	if err = NewCase(path).Unmarshal(b, v); err != nil {
		return
	}
	return
}
