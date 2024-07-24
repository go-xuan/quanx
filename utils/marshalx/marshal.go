package marshalx

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-xuan/quanx/os/filex"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v3"

	"github.com/go-xuan/quanx/types/stringx"
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

var jsonCase, yamlCase, tomlCase, propertiesCase, msgpackCase *Case

func NewCase(name string) *Case {
	switch TypeIdentify(name) {
	case Json:
		return JSON()
	case Yaml, Yml:
		return YAML()
	case Toml:
		return TOML()
	case Properties:
		return PROPERTIES()
	case Msgpack:
		return MSGPACK()
	default:
		return JSON()
	}
}

func JSON() *Case {
	if jsonCase == nil {
		jsonCase = &Case{Name: Json, Marshal: json.Marshal, Unmarshal: json.Unmarshal}
	}
	return jsonCase
}

func YAML() *Case {
	if yamlCase == nil {
		yamlCase = &Case{Name: Yaml, Marshal: yaml.Marshal, Unmarshal: yaml.Unmarshal}
	}
	return yamlCase
}

func TOML() *Case {
	if tomlCase == nil {
		tomlCase = &Case{Name: Toml, Marshal: TomlMarshal, Unmarshal: toml.Unmarshal}
	}
	return tomlCase
}

func PROPERTIES() *Case {
	if propertiesCase == nil {
		propertiesCase = &Case{Name: Properties, Marshal: PropertiesMarshal, Unmarshal: PropertiesUnmarshal}
	}
	return propertiesCase
}

func MSGPACK() *Case {
	if msgpackCase == nil {
		msgpackCase = &Case{Name: Msgpack, Marshal: msgpack.Marshal, Unmarshal: msgpack.Unmarshal}
	}
	return msgpackCase
}

// Case 序列化方案
type Case struct {
	Name      string
	Marshal   func(any) ([]byte, error)
	Unmarshal func([]byte, any) error
}

// TypeIdentify 类型识别
func TypeIdentify(name string) string {
	if stringx.ContainsAny(name, ".", "\\", "/") {
		return filex.Suffix(name)
	}
	return name
}

// UnmarshalFromFile 读取配置文件到指针
func UnmarshalFromFile(path string, v any) (err error) {
	if !filex.Exists(path) {
		return errors.New("the file not exist : " + path)
	}
	var bytes []byte
	if bytes, err = os.ReadFile(path); err != nil {
		return
	}
	if err = NewCase(path).Unmarshal(bytes, v); err != nil {
		return
	}
	return
}

// MarshalToStruct 将任意对象序列化成指定结构体对象
func MarshalToStruct[T any](v any, filename string) (t T, err error) {
	var newCase = NewCase(filename)
	var bytes []byte
	if bytes, err = newCase.Marshal(v); err != nil {
		return
	}
	if err = newCase.Unmarshal(bytes, &t); err != nil {
		return
	}
	return
}
