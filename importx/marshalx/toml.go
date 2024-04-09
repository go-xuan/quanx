package marshalx

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

func TomlMarshal(v any) (out []byte, err error) {
	var buffer bytes.Buffer
	if err = toml.NewEncoder(&buffer).Encode(v); err != nil {
		return
	}
	out = buffer.Bytes()
	return
}
