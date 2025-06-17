package marshalx

import (
	"bytes"

	"github.com/BurntSushi/toml"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
)

func Toml() Marshal {
	return tomlImpl{}
}

type tomlImpl struct{}

func (t tomlImpl) Name() string {
	return tomlMethod
}

func (t tomlImpl) Marshal(v interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	if err := toml.NewEncoder(&buffer).Encode(v); err != nil {
		return nil, errorx.Wrap(err, "encode toml failed")
	}
	return buffer.Bytes(), nil
}

func (t tomlImpl) Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

func (t tomlImpl) Read(path string, v interface{}) error {
	if data, err := readFile(path); err != nil {
		return errorx.Wrap(err, "read file error")
	} else {
		return t.Unmarshal(data, v)
	}
}

func (t tomlImpl) Write(path string, v interface{}) error {
	if data, err := t.Marshal(v); err != nil {
		return errorx.Wrap(err, "toml marshal error")
	} else if err = filex.WriteFile(path, data); err != nil {
		return errorx.Wrap(err, "write file error")
	}
	return nil
}
