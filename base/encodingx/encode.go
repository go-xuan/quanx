package encodingx

// EncodeToString 加密接口
type EncodeToString interface {
	Encode(plaintext []byte) string
}

// DecodeString 解密接口
type DecodeString interface {
	Decode(ciphertext string) ([]byte, error)
}
