package cryptox

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestAES(t *testing.T) {
	key := "fd6d1c1fb333a9ab4e17ce683c988d75"
	iv := "b43f9be55644daca"
	crypto, err := AES(key, iv, AesECB)
	if err != nil {
		fmt.Println(err)
	}

	bytes, _ := json.Marshal(struct {
		AppID     string `json:"app_id"`
		Timestamp int64  `json:"timestamp"`
	}{
		AppID:     "1234567890",
		Timestamp: time.Now().Unix(),
	})

	var ciphertext, plaintext []byte
	if ciphertext, err = crypto.Encrypt(bytes); err != nil {
		panic(err)
	}
	fmt.Println("加密：", Base64().Encode(ciphertext))

	if plaintext, err = crypto.Decrypt(ciphertext); err != nil {
		panic(err)
	}
	fmt.Println("解密：", string(plaintext))
}

func TestRsa(t *testing.T) {
	var priPem = "./rsa/private.pem"
	var publicPem = "./rsa/public.pem"
	crypto, err := RSA(priPem, publicPem, RsaPKCS8, RsaPKIX)
	if err != nil {
		fmt.Println(err)
	}

	bytes, _ := json.Marshal(struct {
		AppID     string `json:"app_id"`
		Timestamp int64  `json:"timestamp"`
	}{
		AppID:     "1234567890",
		Timestamp: time.Now().Unix(),
	})

	var ciphertext, plaintext []byte
	if ciphertext, err = crypto.Encrypt(bytes); err != nil {
		panic(err)
	}
	fmt.Println("加密：", Base64().Encode(ciphertext))

	if plaintext, err = crypto.Decrypt(ciphertext); err != nil {
		panic(err)
	}
	fmt.Println("解密：", string(plaintext))

	err = crypto.Persistence()
	fmt.Println(err)
}
