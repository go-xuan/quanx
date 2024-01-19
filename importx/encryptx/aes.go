package encryptx

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/go-xuan/quanx/utilx/randx"
)

var aesX *Aes

type Aes struct {
	encrypter cipher.BlockMode
	decrypter cipher.BlockMode
	key       []byte
	iv        []byte
}

func AES() *Aes {
	if aesX == nil {
		var err error
		aesX, err = newAES()
		if err != nil {
			return nil
		}
	}
	return aesX
}

func newAES() (myRsa *Aes, err error) {
	var key = []byte(randx.String(16))
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		return
	}
	var iv = []byte(randx.String(block.BlockSize()))
	myRsa = &Aes{
		encrypter: cipher.NewCBCEncrypter(block, iv),
		decrypter: cipher.NewCBCDecrypter(block, iv),
		key:       key,
		iv:        iv,
	}
	return
}

func (m *Aes) Encrypt(plaintext []byte) (ciphertext []byte) {
	m.encrypter.CryptBlocks(ciphertext, plaintext)
	return
}

func (m *Aes) Decrypt(plaintext []byte) (ciphertext []byte) {
	m.decrypter.CryptBlocks(plaintext, ciphertext)
	return
}

func AesEncrypt(plaintext []byte, key []byte, iv []byte) (ciphertext []byte, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		return
	}
	ciphertext = make([]byte, len(plaintext))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, plaintext)
	return
}

func AesDecrypt(ciphertext []byte, key []byte, iv []byte) (plaintext []byte, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		return
	}
	plaintext = make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)
	return
}
