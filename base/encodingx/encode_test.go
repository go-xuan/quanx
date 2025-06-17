package encodingx

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	var plaintext = "123456"
	ciphertext := Base64(true).Encode([]byte(plaintext))
	fmt.Println(ciphertext)

	ciphertext = Hash(md5.New()).Encode([]byte(plaintext))
	fmt.Println(ciphertext)

	ciphertext = Hash(sha1.New()).Encode([]byte(plaintext))
	fmt.Println(ciphertext)

	ciphertext = Hash(sha256.New()).Encode([]byte(plaintext))
	fmt.Println(ciphertext)
	
}
