package marshalx

import (
	"encoding/json"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
)

// Json 序列化方式
// indent: 缩进值，默认不缩进，仅在 Marshal 时有效
func Json(indent ...string) Marshal {
	if len(indent) > 0 {
		return jsonImpl{indent: indent[0]}
	}
	return jsonImpl{}
}

type jsonImpl struct {
	indent string
}

func (j jsonImpl) Name() string {
	return jsonMethod
}

func (j jsonImpl) Marshal(v interface{}) ([]byte, error) {
	if j.indent != "" {
		return json.MarshalIndent(v, "", j.indent)
	}
	return json.Marshal(v)
}

func (j jsonImpl) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (j jsonImpl) Read(path string, v interface{}) error {
	if data, err := readFile(path); err != nil {
		return errorx.Wrap(err, "read file error")
	} else {
		return j.Unmarshal(data, v)
	}
}

// Write 写入json文件
func (j jsonImpl) Write(path string, v interface{}) error {
	if data, err := json.Marshal(v); err != nil {
		return errorx.Wrap(err, "json marshal error")
	} else if err = filex.WriteFile(path, data); err != nil {
		return errorx.Wrap(err, "write file error")
	}
	return nil
}
