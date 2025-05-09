package cryptox

import (
	"encoding/base64"
)

// Base64 加密
// safe: 是否使用URL安全的Base64编码
// 注意：URL安全的Base64编码会将+和/替换为-和_，因此在解码时需要将-和_替换为+和/
func Base64(safe ...bool) Base64Crypto {
	return Base64Crypto{safe: len(safe) > 0 && safe[0]}
}

// Base64Crypto Base64加密器
type Base64Crypto struct {
	safe bool
}

// Encode 加密
func (c Base64Crypto) Encode(plaintext []byte) string {
	if c.safe {
		return base64.URLEncoding.EncodeToString(plaintext)
	} else {
		return base64.StdEncoding.EncodeToString(plaintext)
	}
}

// Decode 解密
func (c Base64Crypto) Decode(ciphertext string) ([]byte, error) {
	if c.safe {
		return base64.URLEncoding.DecodeString(ciphertext)
	} else {
		return base64.StdEncoding.DecodeString(ciphertext)
	}
}
