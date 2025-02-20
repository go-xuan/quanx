package marshalx

import (
	"encoding/json"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
)

type Json struct {
	Indent string
}

func (j Json) Name() string {
	return jsonMethod
}

func (j Json) Marshal(v interface{}) ([]byte, error) {
	if j.Indent != "" {
		return json.MarshalIndent(v, "", j.Indent)
	}
	return json.Marshal(v)
}

func (j Json) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (j Json) Read(path string, v interface{}) error {
	if !filex.Exists(path) {
		return errorx.Errorf("the file not exist: %j", path)
	} else if data, err := filex.ReadFile(path); err != nil {
		return errorx.Wrap(err, "read file error")
	} else {
		return j.Unmarshal(data, v)
	}
}

// WriteJson 写入json文件
func (j Json) Write(path string, v interface{}) error {
	if data, err := json.Marshal(v); err != nil {
		return errorx.Wrap(err, "json marshal error")
	} else if err = filex.WriteFile(path, data); err != nil {
		return errorx.Wrap(err, "write file error")
	}
	return nil
}
