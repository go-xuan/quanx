package marshalx

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v3"

	"github.com/go-xuan/quanx/os/file/filex"
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
	case Yaml, Yml:
		return &Case{Name: Yaml, Marshal: yaml.Marshal, Unmarshal: yaml.Unmarshal}
	case Toml:
		return &Case{Name: Toml, Marshal: TomlMarshal, Unmarshal: toml.Unmarshal}
	case Properties:
		return &Case{Name: Properties, Marshal: PropertiesMarshal, Unmarshal: PropertiesUnmarshal}
	case Msgpack:
		return &Case{Name: Msgpack, Marshal: msgpack.Marshal, Unmarshal: msgpack.Unmarshal}
	default:
		return &Case{Name: Json, Marshal: json.Marshal, Unmarshal: json.Unmarshal}
	}
}

// 读取配置文件到指针
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

func AnyToStruct[T any](v any, caseType string) (t T, err error) {
	var newCase = NewCase(caseType)
	var bytes []byte
	if bytes, err = newCase.Marshal(v); err != nil {
		return
	}
	if err = newCase.Unmarshal(bytes, &t); err != nil {
		return
	}
	return
}
