package marshalx

import (
	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
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

// Marshal 序列化
type Marshal interface {
	Name() string
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	Read(string, interface{}) error
	Write(string, interface{}) error
}

// Apply 适配序列化
func Apply(name string) Marshal {
	if stringx.ContainsAny(name, ".", "\\", "/") {
		name = filex.GetSuffix(name)
	}
	switch name {
	case ymlMethod, yamlMethod:
		return Yaml()
	case tomlMethod:
		return Toml()
	case propertiesMethod:
		return Properties()
	case msgpackMethod:
		return Msgpack()
	case jsonMethod:
		return Json("    ")
	default:
		return Json()
	}
}

func readFile(path string) ([]byte, error) {
	if !filex.Exists(path) {
		return nil, errorx.Errorf("the file not exist: %s", filex.Pwd(path))
	} else if data, err := filex.ReadFile(path); err != nil {
		return nil, errorx.Wrap(err, "read file error")
	} else {
		return data, nil
	}
}
