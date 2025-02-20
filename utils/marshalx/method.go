package marshalx

import (
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/types/stringx"
)

const (
	jsonMethod       = "json"
	yamlMethod       = "yml"
	ymlMethod        = "yaml"
	tomlMethod       = "toml"
	propertiesMethod = "properties"
	msgpackMethod    = "msgpack"
)

// Method 序列化方式
type Method interface {
	Name() string
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	Read(string, interface{}) error
	Write(string, interface{}) error
}

// Apply 适配序列化方式
func Apply(name string) Method {
	if stringx.ContainsAny(name, ".", "\\", "/") {
		name = filex.GetSuffix(name)
	}
	switch name {
	case jsonMethod:
		return Json{Indent: "    "}
	case ymlMethod, yamlMethod:
		return Yaml{}
	case tomlMethod:
		return Toml{}
	case propertiesMethod:
		return Properties{}
	case msgpackMethod:
		return Msgpack{}
	default:
		return Json{}
	}
}
