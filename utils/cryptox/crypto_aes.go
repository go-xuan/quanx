package cryptox

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/types/anyx"
)

// AES 创建AES对象（默认CBC模式）
func AES(key, iv string, mode ...Mode) (*AesCrypto, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, errorx.Wrap(err, "new cipher error")
	}
	if ivl, bs := len(iv), block.BlockSize(); ivl != bs {
		return nil, errorx.Errorf("iv length must equal block size: %d != %d", ivl, bs)
	}
	return &AesCrypto{
		mode:  anyx.Default(AesCBC, mode...),
		key:   key,
		iv:    []byte(iv),
		block: block,
	}, nil
}

// AesCrypto AES加密器
type AesCrypto struct {
	mode  Mode         // 加密模式
	key   string       // 秘钥, 长度必须为16, 24或32字节
	iv    []byte       // 初始化向量, 长度必须和加密块大小相等
	block cipher.Block // 加密块
}

// Encrypt 加密
func (c *AesCrypto) Encrypt(plaintext []byte) ([]byte, error) {
	plaintext = pkcs7Padding(plaintext, c.block.BlockSize())
	var ciphertext = make([]byte, len(plaintext))
	switch c.mode {
	case AesCBC:
		cipher.NewCBCEncrypter(c.block, c.iv).CryptBlocks(ciphertext, plaintext)
	case AesCFB:
		cipher.NewCFBEncrypter(c.block, c.iv).XORKeyStream(ciphertext, plaintext)
	case AesECB:
		size, blockSize := len(plaintext), c.block.BlockSize()
		for i := 0; i < size; i += blockSize {
			c.block.Encrypt(ciphertext[i:i+blockSize], plaintext[i:i+blockSize])
		}
	default:
		return nil, errorx.New("unsupported aes mode")
	}
	return ciphertext, nil
}

// Decrypt 解密
func (c *AesCrypto) Decrypt(ciphertext []byte) ([]byte, error) {
	size, blockSize := len(ciphertext), c.block.BlockSize()
	if size%blockSize != 0 {
		return nil, errorx.Errorf("the ciphertext size error: %d/%d", size, blockSize)
	}
	var plaintext = make([]byte, size)
	switch c.mode {
	case AesCBC:
		cipher.NewCBCDecrypter(c.block, c.iv).CryptBlocks(plaintext, ciphertext)
	case AesCFB:
		cipher.NewCFBDecrypter(c.block, c.iv).XORKeyStream(plaintext, ciphertext)
	case AesECB:
		for i := 0; i < size; i += blockSize {
			c.block.Decrypt(plaintext[i:i+blockSize], ciphertext[i:i+blockSize])
		}
	default:
		return nil, errorx.New("unsupported aes mode")
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
