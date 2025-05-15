package cryptox

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-xuan/quanx/base/encodingx"
	"github.com/go-xuan/quanx/base/filex"
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
	fmt.Println("加密：", encodingx.Base64().Encode(ciphertext))

	if plaintext, err = crypto.Decrypt(ciphertext); err != nil {
		panic(err)
	}
	fmt.Println("解密：", string(plaintext))
}

func TestRsa(t *testing.T) {
	data, err := filex.ReadFile("./rsa/private.pem")
	if err != nil {
		fmt.Println(err)
	}
	crypto, err := ParseRsaCrypto(data, RsaPKCS8, RsaPKCS8)
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
	fmt.Println("加密：", encodingx.Base64().Encode(ciphertext))

	if plaintext, err = crypto.Decrypt(ciphertext); err != nil {
		panic(err)
	}
	fmt.Println("解密：", string(plaintext))

	if err = crypto.SavePrivateKey("./rsa/private.pem"); err != nil {
		fmt.Println(err)
	}
	if err = crypto.SavePublicKey("./rsa/public-1.pem"); err != nil {
		fmt.Println(err)
	}
}
