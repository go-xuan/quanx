package cryptox

import (
	"encoding/hex"
	"hash"
)

// Hash 哈希加密
func Hash(h hash.Hash) HashCrypto {
	return HashCrypto{h: h}
}

// HashCrypto 哈希加密器
type HashCrypto struct {
	h hash.Hash
}

// Encrypt 加密
func (c HashCrypto) Encrypt(plaintext []byte) string {
	c.h.Write(plaintext)
	return hex.EncodeToString(c.h.Sum(nil))
}
