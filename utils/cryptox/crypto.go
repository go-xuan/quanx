package cryptox

type Mode uint

const (
	AesCBC Mode = iota + 1
	AesCFB
	AesECB

	RsaPKCS1 // 私钥和公钥都能使用
	RsaPKCS8 // 仅能用于私钥
	RsaPKIX  // 仅能用于公钥
)

// Crypto 加解密接口
type Crypto interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}
