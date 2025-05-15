package encodingx

import (
	"encoding/base64"
)

// Base64 加密
// safe: 是否使用URL安全的Base64编码
func Base64(safe ...bool) *Base64Encode {
	return &Base64Encode{safe: len(safe) > 0 && safe[0]}
}

// Base64Encode Base64加密器
type Base64Encode struct {
	safe bool
}

func (c *Base64Encode) Encode(plaintext []byte) string {
	if c.safe {
		return base64.URLEncoding.EncodeToString(plaintext)
	} else {
		return base64.StdEncoding.EncodeToString(plaintext)
	}
}

func (c *Base64Encode) Decode(ciphertext string) ([]byte, error) {
	if c.safe {
		return base64.URLEncoding.DecodeString(ciphertext)
	} else {
		return base64.StdEncoding.DecodeString(ciphertext)
	}
}
