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
	mode  mode         // 加密模式
	key   string       // 秘钥
	iv    string       // 初始化向量
	block cipher.Block // 加密块
}

func AES() *Aes {
	if _aes == nil {
		if _, err := NewAes(randx.String(aes.BlockSize), randx.String(aes.BlockSize)); err != nil {
			panic(err)
		}
	}
	return _aes
}

func NewAes(key, iv string) (*Aes, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, errorx.Wrap(err, "new cipher error")
	}
	if size := block.BlockSize(); size != len(iv) {
		return nil, errorx.Errorf("iv length must equal to the block size: %d", size)
	}
	if _aes == nil {
		_aes = &Aes{
			mode:  CBC,
			key:   key,
			iv:    iv,
			block: block}
	} else {
		_aes.key = key
		_aes.iv = iv
		_aes.block = block
	}
	return _aes, nil
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
	cipher.NewCBCEncrypter(a.block, []byte(a.iv)).CryptBlocks(ciphertext, plaintext)
	return ciphertext
}

// decryptCBC CBC解密
func (a *Aes) decryptCBC(ciphertext []byte) []byte {
	var plaintext = make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(a.block, []byte(a.iv)).CryptBlocks(plaintext, ciphertext)
	return plaintext
}

// encryptCFB CBC加密
func (a *Aes) encryptCFB(plaintext []byte) []byte {
	var ciphertext = make([]byte, len(plaintext))
	cipher.NewCFBEncrypter(a.block, []byte(a.iv)).XORKeyStream(ciphertext, plaintext)
	return ciphertext
}

// decryptCFB CBC解密
func (a *Aes) decryptCFB(ciphertext []byte) []byte {
	var plaintext = make([]byte, len(ciphertext))
	cipher.NewCFBDecrypter(a.block, []byte(a.iv)).XORKeyStream(plaintext, ciphertext)
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
