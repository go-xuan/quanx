package encryptx

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestAES(t *testing.T) {
	aes, err := NewAes("fd6d1c1fb333a9ab4e17ce683c988d75", "b43f9be55644daca")
	if err != nil {
		fmt.Println(err)
	}
	type SignData struct {
		AppID     string `json:"app_id"`
		SecretKey string `json:"secret_key"`
		Timestamp int64  `json:"timestamp"`
	}
	bytes, _ := json.Marshal(SignData{
		AppID:     "e9e97231-d009-4d39-b50f-bc8fc9008300",
		SecretKey: "fd6d1c1fb333a9ab4e17ce683c988d75",
		Timestamp: time.Now().Unix(),
	})
	// 加密
	ciphertext := aes.Encrypt(bytes)
	fmt.Println(Base64Encode(ciphertext))

	// 解密
	plaintext := aes.Decrypt(ciphertext)
	fmt.Println(string(plaintext))
}
