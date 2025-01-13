package marshalx

import (
	"bytes"

	"github.com/BurntSushi/toml"
	
	"github.com/go-xuan/quanx/os/errorx"
)

func TomlMarshal(v any) ([]byte, error) {
	var buffer bytes.Buffer
	if err := toml.NewEncoder(&buffer).Encode(v); err != nil {
		return nil, errorx.Wrap(err, "encode toml failed")
	}
	return buffer.Bytes(), nil
}
