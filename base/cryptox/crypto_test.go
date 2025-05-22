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
	aesCrypto, err := NewAesCrypto(key, iv, GCM)
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
	if ciphertext, err = aesCrypto.Encrypt(bytes); err != nil {
		panic(err)
	}
	fmt.Println("加密：", encodingx.Base64().Encode(ciphertext))

	if plaintext, err = aesCrypto.Decrypt(ciphertext); err != nil {
		panic(err)
	}
	fmt.Println("解密：", string(plaintext))
}

func TestRsa(t *testing.T) {
	data, err := filex.ReadFile("./rsa/private.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	rsaCrypto, err := ParseRsaCrypto(data, PKCS8)
	//rsaCrypto, err := NewRsaCrypto(2048, PKCS8, PKIX)
	if err != nil {
		fmt.Println(err)
		return
	}
	bytes, _ := json.Marshal(struct {
		AppID     string `json:"app_id"`
		Timestamp int64  `json:"timestamp"`
	}{
		AppID:     "1234567890",
		Timestamp: time.Now().Unix(),
	})

	var ciphertext, plaintext []byte
	if ciphertext, err = rsaCrypto.Encrypt(bytes); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("加密：", encodingx.Base64().Encode(ciphertext))

	if plaintext, err = rsaCrypto.Decrypt(ciphertext); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("解密：", string(plaintext))

	if err = rsaCrypto.SavePrivateKey("./rsa/private.pem", PKCS8); err != nil {
		fmt.Println(err)
	}
	if err = rsaCrypto.SavePublicKey("./rsa/public-1.pem", PKCS1); err != nil {
		fmt.Println(err)
	}
}
