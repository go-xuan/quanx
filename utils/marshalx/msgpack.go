package marshalx

import (
	"github.com/vmihailenco/msgpack"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
)

type msgpackImpl struct{}

func (m msgpackImpl) Name() string {
	return msgpackMethod
}

func (m msgpackImpl) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (m msgpackImpl) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

func (m msgpackImpl) Read(path string, v interface{}) error {
	if !filex.Exists(path) {
		return errorx.Errorf("the file not exist: %s", filex.Pwd(path))
	} else if data, err := filex.ReadFile(path); err != nil {
		return errorx.Wrap(err, "read file error")
	} else {
		return m.Unmarshal(data, v)
	}
}

func (m msgpackImpl) Write(path string, v interface{}) error {
	if data, err := m.Marshal(v); err != nil {
		return errorx.Wrap(err, "msgpack marshal error")
	} else if err = filex.WriteFile(path, data); err != nil {
		return errorx.Wrap(err, "write file error")
	}
	return nil
}
