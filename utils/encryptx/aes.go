package encryptx

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/go-xuan/quanx/utils/randx"
)

type mode int

const (
	CBC mode = iota
	CFB
	ECB
)

var _aes *Aes

type Aes struct {
	key   []byte       // 秘钥
	iv    []byte       // 初始化向量
	mode  mode         // 加密模式
	block cipher.Block // 加密块
}

func AES() *Aes {
	if _aes == nil {
		var err error
		if _aes, err = NewAes(); err != nil {
			panic(err)
		}
	}
	return _aes
}

func NewAes(secretKey ...string) (*Aes, error) {
	var key []byte
	if len(secretKey) > 0 {
		key = []byte(secretKey[0])
	} else {
		key = []byte(randx.String(aes.BlockSize))
	}
	if block, err := aes.NewCipher(key); err != nil {
		return nil, err
	} else {
		var iv = []byte(randx.String(block.BlockSize()))
		return &Aes{mode: CBC, block: block, key: key, iv: iv}, nil
	}
}

// Mode 加密模式
func (a *Aes) Mode(m mode) *Aes {
	a.mode = m
	return a
}

// Encrypt 加密
func (a *Aes) Encrypt(plaintext []byte) []byte {
	switch a.mode {
	case CBC:
		return a.EncryptCBC(plaintext)
	case CFB:
		return a.EncryptCFB(plaintext)
	case ECB:
		return a.EncryptECB(plaintext)
	}
	return nil
}

// Decrypt 解密
func (a *Aes) Decrypt(ciphertext []byte) []byte {
	switch a.mode {
	case CBC:
		return a.DecryptCBC(ciphertext)
	case CFB:
		return a.DecryptCFB(ciphertext)
	case ECB:
		return a.DecryptECB(ciphertext)
	}
	return nil
}

// EncryptCBC CBC加密
func (a *Aes) EncryptCBC(plaintext []byte) []byte {
	var ciphertext []byte
	cipher.NewCBCEncrypter(a.block, a.iv).CryptBlocks(ciphertext, plaintext)
	return ciphertext
}

// DecryptCBC CBC解密
func (a *Aes) DecryptCBC(ciphertext []byte) []byte {
	var plaintext []byte
	cipher.NewCBCDecrypter(a.block, a.iv).CryptBlocks(plaintext, ciphertext)
	return plaintext
}

// EncryptCFB CBC加密
func (a *Aes) EncryptCFB(plaintext []byte) []byte {
	var ciphertext []byte
	cipher.NewCFBEncrypter(a.block, a.iv).XORKeyStream(ciphertext, plaintext)
	return ciphertext
}

// DecryptCFB CBC解密
func (a *Aes) DecryptCFB(ciphertext []byte) []byte {
	var plaintext []byte
	cipher.NewCFBDecrypter(a.block, a.iv).XORKeyStream(plaintext, ciphertext)
	return plaintext
}

// EncryptECB ECB加密
func (a *Aes) EncryptECB(plaintext []byte) []byte {
	var ciphertext []byte
	a.block.Encrypt(ciphertext, plaintext)
	return ciphertext
}

// DecryptECB ECB解密
func (a *Aes) DecryptECB(ciphertext []byte) []byte {
	var plaintext []byte
	a.block.Encrypt(plaintext, ciphertext)
	return plaintext
}
