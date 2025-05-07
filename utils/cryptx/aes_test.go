package cryptx

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
	bytes, _ := json.Marshal(struct {
		AppID     string `json:"app_id"`
		Timestamp int64  `json:"timestamp"`
	}{
		AppID:     "1234567890",
		Timestamp: time.Now().Unix(),
	})
	// 加密
	ciphertext, _ := aes.Encrypt(bytes)
	fmt.Println("加密：", Base64Encode(ciphertext))

	// 解密
	plaintext, _ := aes.Decrypt(ciphertext)
	fmt.Println("解密：", string(plaintext))
}
