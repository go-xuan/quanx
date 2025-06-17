package cryptox

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/go-xuan/quanx/utils/randx"

	"github.com/go-xuan/quanx/base/errorx"
)

const (
	CBC Mode = iota + 1
	CFB
	ECB
	GCM
)

// AES 创建AES对象（默认CBC模式）
func AES() (Crypto, error) {
	key := randx.String(16)
	iv := randx.String(16)
	crypto, err := NewAesCrypto(key, iv, CBC)
	if err != nil {
		return nil, errorx.Wrap(err, "new aes crypto error")
	}
	return crypto, nil
}

// NewAesCrypto 创建AES对象（默认CBC模式）
func NewAesCrypto(key, iv string, mode Mode) (Crypto, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, errorx.Wrap(err, "new cipher error")
	}
	switch mode {
	case CBC:
		return newAesCBC(key, iv, block)
	case CFB:
		return newAesCFB(key, iv, block)
	case ECB:
		return newAesECB(key, iv, block)
	case GCM:
		return newAesGCM(key, iv, block)
	default:
		return nil, errorx.New(fmt.Sprintf("unsupported aes mode: %d", mode))
	}
}

func newAesGCM(key, nonce string, block cipher.Block) (*AesGCM, error) {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errorx.Wrap(err, "new gcm error")
	}

	if len(nonce) < gcm.NonceSize() {
		return nil, errorx.New("nonce length must greater than gcm nonce size")
	}
	nonce = nonce[:gcm.NonceSize()]
	return &AesGCM{key: key, nonce: []byte(nonce), block: block, gcm: gcm}, nil
}

type AesGCM struct {
	key   string       // 秘钥, 长度必须为16, 24或32字节
	nonce []byte       // 随机值
	block cipher.Block // 加密块
	gcm   cipher.AEAD  // GCM加密器
}

func (c *AesGCM) Encrypt(plaintext []byte) ([]byte, error) {
	ciphertext := c.gcm.Seal(nil, c.nonce, plaintext, nil)
	return ciphertext, nil
}

func (c *AesGCM) Decrypt(ciphertext []byte) ([]byte, error) {
	plaintext, err := c.gcm.Open(nil, c.nonce, ciphertext, nil)
	if err != nil {
		return nil, errorx.Wrap(err, "gcm open error")
	}
	return plaintext, nil
}

func newAesCBC(key, iv string, block cipher.Block) (*AesCBC, error) {
	if len(iv) != block.BlockSize() {
		return nil, errorx.New("iv length must equal block size")
	}
	return &AesCBC{key: key, iv: []byte(iv), block: block}, nil
}

type AesCBC struct {
	key   string       // 秘钥, 长度必须为16, 24或32字节
	iv    []byte       // 随机值
	block cipher.Block // 加密块
}

func (c *AesCBC) Encrypt(plaintext []byte) ([]byte, error) {
	plaintext = pkcs7Padding(plaintext, c.block.BlockSize())
	var ciphertext = make([]byte, len(plaintext))
	cipher.NewCBCEncrypter(c.block, c.iv).CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func (c *AesCBC) Decrypt(ciphertext []byte) ([]byte, error) {
	if size, blockSize := len(ciphertext), c.block.BlockSize(); size%blockSize != 0 {
		return nil, errorx.Errorf("the ciphertext size error: %d/%d", size, blockSize)
	}
	var plaintext = make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(c.block, c.iv).CryptBlocks(plaintext, ciphertext)
	return pkcs7UnPadding(plaintext)
}

func newAesCFB(key, iv string, block cipher.Block) (*AesCFB, error) {
	if len(iv) != block.BlockSize() {
		return nil, errorx.New("iv length must equal block size")
	}
	return &AesCFB{key: key, iv: []byte(iv), block: block}, nil
}

type AesCFB struct {
	key   string       // 秘钥, 长度必须为16, 24或32字节
	iv    []byte       // 随机值
	block cipher.Block // 加密块
}

func (c *AesCFB) Encrypt(plaintext []byte) ([]byte, error) {
	plaintext = pkcs7Padding(plaintext, c.block.BlockSize())
	var ciphertext = make([]byte, len(plaintext))
	cipher.NewCFBEncrypter(c.block, c.iv).XORKeyStream(ciphertext, plaintext)
	return ciphertext, nil
}

func (c *AesCFB) Decrypt(ciphertext []byte) ([]byte, error) {
	if size, blockSize := len(ciphertext), c.block.BlockSize(); size%blockSize != 0 {
		return nil, errorx.Errorf("the ciphertext size error: %d/%d", size, blockSize)
	}
	var plaintext = make([]byte, len(ciphertext))
	cipher.NewCFBDecrypter(c.block, c.iv).XORKeyStream(plaintext, ciphertext)
	return pkcs7UnPadding(plaintext)
}

func newAesECB(key, iv string, block cipher.Block) (*AesECB, error) {
	if len(iv) != block.BlockSize() {
		return nil, errorx.New("iv length must equal block size")
	}
	return &AesECB{key: key, iv: []byte(iv), block: block}, nil
}

type AesECB struct {
	key   string       // 秘钥, 长度必须为16, 24或32字节
	iv    []byte       // 随机值
	block cipher.Block // 加密块
}

func (c *AesECB) Encrypt(plaintext []byte) ([]byte, error) {
	plaintext = pkcs7Padding(plaintext, c.block.BlockSize())
	size, blockSize := len(plaintext), c.block.BlockSize()
	var ciphertext = make([]byte, size)
	for i := 0; i < size; i += blockSize {
		c.block.Encrypt(ciphertext[i:i+blockSize], plaintext[i:i+blockSize])
	}
	return ciphertext, nil
}

func (c *AesECB) Decrypt(ciphertext []byte) ([]byte, error) {
	size, blockSize := len(ciphertext), c.block.BlockSize()
	if size%blockSize != 0 {
		return nil, errorx.Errorf("the ciphertext size error: %d/%d", size, blockSize)
	}
	var plaintext = make([]byte, size)
	for i := 0; i < size; i += blockSize {
		c.block.Decrypt(plaintext[i:i+blockSize], ciphertext[i:i+blockSize])
	}
	return pkcs7UnPadding(plaintext)
}

// pkcs7Padding 使用 PKCS7 进行填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, paddingText...)
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
