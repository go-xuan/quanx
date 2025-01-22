package marshalx

import (
	"github.com/vmihailenco/msgpack"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
)

type Msgpack struct{}

func (s Msgpack) Name() string {
	return msgpackStrategy
}

func (s Msgpack) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (s Msgpack) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

func (s Msgpack) Read(path string, v interface{}) error {
	if !filex.Exists(path) {
		return errorx.Errorf("the file not exist: %s", path)
	} else if data, err := filex.ReadFile(path); err != nil {
		return errorx.Wrap(err, "read file error")
	} else {
		return s.Unmarshal(data, v)
	}
}

func (s Msgpack) Write(path string, v interface{}) error {
	if data, err := s.Marshal(v); err != nil {
		return errorx.Wrap(err, "msgpack marshal error")
	} else if err = filex.WriteFile(path, data); err != nil {
		return errorx.Wrap(err, "write file error")
	}
	return nil
}
