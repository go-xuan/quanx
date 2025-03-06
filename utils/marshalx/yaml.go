package marshalx

import (
	"gopkg.in/yaml.v3"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
)

type yamlImpl struct{}

func (y yamlImpl) Name() string {
	return yamlMethod
}

func (y yamlImpl) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y yamlImpl) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

// WriteYaml 写入yaml文件
func WriteYaml(path string, v any) error {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return errorx.Wrap(err, "yamlMethod marshal error")
	}
	if err = filex.WriteFile(path, bytes); err != nil {
		return errorx.Wrap(err, "write file error")
	}
	return nil
}

func (y yamlImpl) Read(path string, v interface{}) error {
	if !filex.Exists(path) {
		return errorx.Errorf("the file not exist: %s", filex.Pwd(path))
	} else if data, err := filex.ReadFile(path); err != nil {
		return errorx.Wrap(err, "read file error")
	} else {
		return y.Unmarshal(data, v)
	}
}

func (y yamlImpl) Write(path string, v interface{}) error {
	if data, err := y.Marshal(v); err != nil {
		return errorx.Wrap(err, "yaml marshal error")
	} else if err = filex.WriteFile(path, data); err != nil {
		return errorx.Wrap(err, "write file error")
	}
	return nil
}
