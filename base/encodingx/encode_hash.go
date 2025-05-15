package encodingx

import (
	"encoding/hex"
	"hash"
)

// Hash 哈希加密
func Hash(h hash.Hash) *HashEncode {
	return &HashEncode{h: h}
}

// HashEncode 哈希加密器
type HashEncode struct {
	h hash.Hash
}

func (c *HashEncode) Encode(plaintext []byte) string {
	c.h.Write(plaintext)
	return hex.EncodeToString(c.h.Sum(nil))
}
