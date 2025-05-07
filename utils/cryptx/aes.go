package cryptx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/go-xuan/quanx/base/errorx"
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
	key   string       // 秘钥, 长度必须为16, 24或32字节
	iv    []byte       // 初始化向量, 长度必须和加密块大小相等
	block cipher.Block // 加密块
}

func AES() *Aes {
	if _aes == nil {
		if newAes, err := NewAes(randx.String(aes.BlockSize), randx.String(aes.BlockSize)); err != nil {
			panic(err)
		} else {
			_aes = newAes
		}
	}
	return _aes
}

// NewAes 创建AES对象（默认CBC模式）
func NewAes(key, iv string) (*Aes, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, errorx.Wrap(err, "new cipher error")
	}
	if ivl, bs := len(iv), block.BlockSize(); ivl != bs {
		return nil, errorx.Errorf("iv length must equal block size: %d != %d", ivl, bs)
	}
	return &Aes{mode: CBC, key: key, iv: []byte(iv), block: block}, nil
}

// Mode 加密模式
func (a *Aes) Mode(m mode) *Aes {
	a.mode = m
	return a
}

// Encrypt 加密
func (a *Aes) Encrypt(plaintext []byte) ([]byte, error) {
	plaintext = pkcs7Padding(plaintext, a.block.BlockSize())
	var ciphertext = make([]byte, len(plaintext))
	switch a.mode {
	case CBC:
		cipher.NewCBCEncrypter(a.block, a.iv).CryptBlocks(ciphertext, plaintext)
	case CFB:
		cipher.NewCFBEncrypter(a.block, a.iv).XORKeyStream(ciphertext, plaintext)
	case ECB:
		size, blockSize := len(plaintext), a.block.BlockSize()
		for i := 0; i < size; i += blockSize {
			a.block.Encrypt(ciphertext[i:i+blockSize], plaintext[i:i+blockSize])
		}
	default:
		return nil, errorx.New("unsupported mode")
	}
	return ciphertext, nil
}

// Decrypt 解密
func (a *Aes) Decrypt(ciphertext []byte) ([]byte, error) {
	size, blockSize := len(ciphertext), a.block.BlockSize()
	if size%blockSize != 0 {
		return nil, errorx.Errorf("the ciphertext size error: %d/%d", size, blockSize)
	}
	var plaintext = make([]byte, size)
	switch a.mode {
	case CBC:
		cipher.NewCBCDecrypter(a.block, a.iv).CryptBlocks(plaintext, ciphertext)
	case CFB:
		cipher.NewCFBDecrypter(a.block, a.iv).XORKeyStream(plaintext, ciphertext)
	case ECB:
		for i := 0; i < size; i += blockSize {
			a.block.Decrypt(plaintext[i:i+blockSize], ciphertext[i:i+blockSize])
		}
	default:
		return nil, errorx.New("unsupported mode")
	}
	return pkcs7UnPadding(plaintext)
}

// pkcs7Padding 使用 PKCS7 进行填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 移除 PKCS7 填充
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("empty data")
	}
	padding := int(data[length-1])
	if padding > length || padding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:length-padding], nil
}
