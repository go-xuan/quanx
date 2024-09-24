package encryptx

import (
	"fmt"
	"testing"
)

func TestAES(t *testing.T) {
	var key, iv []byte
	key = []byte("cbdb7271281e2ab3")
	iv = []byte("b43f9be55644daca")
	aes, err := NewAes(key, iv)
	if err != nil {
		fmt.Println(err)
	}

	ciphertext := aes.Encrypt([]byte("admin"))
	fmt.Println(ciphertext)
	fmt.Println(string(ciphertext))

	plaintext := aes.Decrypt(ciphertext)
	fmt.Println(plaintext)
	fmt.Println(string(plaintext))
}
