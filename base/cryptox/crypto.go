package cryptox

type Mode uint

// Crypto 加解密接口
type Crypto interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}
