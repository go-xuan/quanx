package marshalx

import (
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/types/stringx"
)

// 序列化器类型
const (
	jsonStrategy       = "json"
	yamlStrategy       = "yml"
	ymlStrategy        = "yaml"
	tomlStrategy       = "toml"
	propertiesStrategy = "properties"
	msgpackStrategy    = "msgpack"
)

// Strategy 序列化策略
type Strategy interface {
	Name() string
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	Read(string, interface{}) error
	Write(string, interface{}) error
}

// Apply 适用序列化策略
func Apply(name string) Strategy {
	if stringx.ContainsAny(name, ".", "\\", "/") {
		name = filex.GetSuffix(name)
	}
	switch name {
	case jsonStrategy:
		return Json{}
	case ymlStrategy, yamlStrategy:
		return Yaml{}
	case tomlStrategy:
		return Toml{}
	case propertiesStrategy:
		return Properties{}
	case msgpackStrategy:
		return Msgpack{}
	default:
		return Json{}
	}
}
