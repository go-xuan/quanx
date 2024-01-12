package marshalx

import (
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/go-xuan/quanx/utilx/filex"
	"github.com/go-xuan/quanx/utilx/marshalx/propertiesx"
	"gopkg.in/yaml.v3"
	"os"
)

// 读取配置文件到指针
func LoadFromFile(config interface{}, filePath string) (err error) {
	if !filex.Exists(filePath) {
		return errors.New("the file not exist : " + filePath)
	}
	var bytes []byte
	bytes, err = os.ReadFile(filePath)
	if err != nil {
		return
	}
	var suffix = filex.Suffix(filePath)
	if err = UnmarshalToPointer(config, bytes, suffix); err != nil {
		return
	}
	return
}

// 解析bytes到指针
func UnmarshalToPointer(config interface{}, bytes []byte, suffix string) (err error) {
	switch suffix {
	case filex.Json:
		err = json.Unmarshal(bytes, config)
	case filex.Yaml, filex.Yml:
		err = yaml.Unmarshal(bytes, config)
	case filex.Toml:
		err = toml.Unmarshal(bytes, config)
	case filex.Properties:
		err = propertiesx.Unmarshal(bytes, config)
	default:
		err = errors.New("the file type is not supported :" + suffix)
	}
	if err != nil {
		return
	}
	return
}
