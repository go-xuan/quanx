package encryptx

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/go-xuan/quanx/os/errorx"
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
		var key = []byte(randx.String(aes.BlockSize))
		var iv = []byte(randx.String(aes.BlockSize))
		if _aes, err = NewAes(key, iv); err != nil {
			panic(err)
		}
	}
	return _aes
}

func NewAes(key, iv []byte) (*Aes, error) {
	if block, err := aes.NewCipher(key); err != nil {
		return nil, errorx.Wrap(err, "new cipher error")
	} else if len(iv) != aes.BlockSize {
		return nil, errorx.New("iv length must be equal to the block size")
	} else {
		return &Aes{mode: CBC, block: block, key: key, iv: iv}, nil
	}
}

func addPadding(text []byte) []byte {
	if len(text)%16 == 0 {
		return text
	}
	padding := 16 - (len(text) % 16)
	for i := 0; i < padding; i++ {
		text = append(text, byte(padding))
	}
	return text
}

// Mode 加密模式
func (a *Aes) Mode(m mode) *Aes {
	a.mode = m
	return a
}

// Encrypt 加密
func (a *Aes) Encrypt(plaintext []byte) []byte {
	plaintext = addPadding(plaintext)
	switch a.mode {
	case CBC:
		return a.encryptCBC(plaintext)
	case CFB:
		return a.encryptCFB(plaintext)
	case ECB:
		return a.encryptECB(plaintext)
	}
	return nil
}

// Decrypt 解密
func (a *Aes) Decrypt(ciphertext []byte) []byte {
	switch a.mode {
	case CBC:
		return a.decryptCBC(ciphertext)
	case CFB:
		return a.decryptCFB(ciphertext)
	case ECB:
		return a.decryptECB(ciphertext)
	}
	return nil
}

// encryptCBC CBC加密
func (a *Aes) encryptCBC(plaintext []byte) []byte {
	var ciphertext = make([]byte, len(plaintext))
	cipher.NewCBCEncrypter(a.block, a.iv).CryptBlocks(ciphertext, plaintext)
	return ciphertext
}

// decryptCBC CBC解密
func (a *Aes) decryptCBC(ciphertext []byte) []byte {
	var plaintext = make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(a.block, a.iv).CryptBlocks(plaintext, ciphertext)
	return plaintext
}

// encryptCFB CBC加密
func (a *Aes) encryptCFB(plaintext []byte) []byte {
	var ciphertext = make([]byte, len(plaintext))
	cipher.NewCFBEncrypter(a.block, a.iv).XORKeyStream(ciphertext, plaintext)
	return ciphertext
}

// decryptCFB CBC解密
func (a *Aes) decryptCFB(ciphertext []byte) []byte {
	var plaintext = make([]byte, len(ciphertext))
	cipher.NewCFBDecrypter(a.block, a.iv).XORKeyStream(plaintext, ciphertext)
	return plaintext
}

// encryptECB ECB加密
func (a *Aes) encryptECB(plaintext []byte) []byte {
	var ciphertext = make([]byte, len(plaintext))
	a.block.Encrypt(ciphertext, plaintext)
	return ciphertext
}

// decryptECB ECB解密
func (a *Aes) decryptECB(ciphertext []byte) []byte {
	var plaintext = make([]byte, len(ciphertext))
	a.block.Encrypt(plaintext, ciphertext)
	return plaintext
}
