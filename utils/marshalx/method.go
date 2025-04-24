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

func Json(indent ...string) Method {
	if len(indent) > 0 {
		return jsonImpl{indent: indent[0]}
	}
	return jsonImpl{}
}

func Yaml() Method {
	return yamlImpl{}
}

func Toml() Method {
	return tomlImpl{}
}

func Properties() Method {
	return propertiesImpl{}
}

func Msgpack() Method {
	return msgpackImpl{}
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
