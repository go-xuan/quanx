package marshalx

import (
	"github.com/vmihailenco/msgpack"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
)

type Msgpack struct{}

func (m Msgpack) Name() string {
	return msgpackMethod
}

func (m Msgpack) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (m Msgpack) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

func (m Msgpack) Read(path string, v interface{}) error {
	if !filex.Exists(path) {
		return errorx.Errorf("the file not exist: %m", path)
	} else if data, err := filex.ReadFile(path); err != nil {
		return errorx.Wrap(err, "read file error")
	} else {
		return m.Unmarshal(data, v)
	}
}

func (m Msgpack) Write(path string, v interface{}) error {
	if data, err := m.Marshal(v); err != nil {
		return errorx.Wrap(err, "msgpack marshal error")
	} else if err = filex.WriteFile(path, data); err != nil {
		return errorx.Wrap(err, "write file error")
	}
	return nil
}
