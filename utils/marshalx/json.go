package marshalx

import (
	"encoding/json"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
)

type Json struct{}

func (s Json) Name() string {
	return jsonStrategy
}

func (s Json) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (s Json) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (s Json) Read(path string, v interface{}) error {
	if !filex.Exists(path) {
		return errorx.Errorf("the file not exist: %s", path)
	} else if data, err := filex.ReadFile(path); err != nil {
		return errorx.Wrap(err, "read file error")
	} else {
		return s.Unmarshal(data, v)
	}
}

// WriteJson 写入json文件
func (s Json) Write(path string, v interface{}) error {
	if data, err := json.MarshalIndent(v, "", "	"); err != nil {
		return errorx.Wrap(err, "json marshal error")
	} else if err = filex.WriteFile(path, data); err != nil {
		return errorx.Wrap(err, "write file error")
	}
	return nil
}
