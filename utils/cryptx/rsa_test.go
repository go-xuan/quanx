package cryptx

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/farmerx/gorsa"
)

// 加密
func TestEncrypt(t *testing.T) {
	var password = "user_0523"
	var pemPath = "./pem/rsa-public.pem"
	var err error
	var ciphertext, plaintext []byte
	plaintext = []byte(password)
	if ciphertext, err = RsaEncryptPKIX(plaintext, pemPath); err != nil {
		panic(err)
	}
	fmt.Println(Base64Encode(ciphertext))
}

func TestDecrypt(t *testing.T) {
	var password = "LAUakl71YSmc1iRz2MdtUplyLFztMpO6hPJz2v0YwknGIgYWVN+FMOz2/hfy8Gwjo3x+8R21/dwbaM+nD6h5lTwrS+/qmIvwd5HyrBZpLz8hMa27OfsIgtccfUI4crt8Oj7qnAKtawmx5BCi9Iyp6uuNri9CqEiImxrKtXIJuVeTrquaC+WIoU7ugvDMN3qoun8uUYZkMLfCRgQ29DOeQbh43jFGiCtQ1v3DaNbnLsHsWa88hX0bbNp29pQph67dB9BvkYHfEiGimulkxYT7uDJHUth4XSJmIG7L+Mb8dvD2oFNIklJwTMLDNhZ3QrMdWYJqoNVuzcuBt1yjHdOe+A=="
	var pemPath = "./pem/rsa-private.pem"
	if plaintext, err := DecryptPassword(password, pemPath); err == nil {
		fmt.Println(Base64Encode(plaintext))
	} else {
		panic(err)
	}
}

// 密码RSA解密获取明文
func DecryptPassword(password string, pemPath string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return nil, err
	}

	var privateKey []byte
	if privateKey, err = os.ReadFile(pemPath); err != nil {
		return nil, err
	}
	var plaintext []byte
	if plaintext, err = PriKeyDecrypt(cipherText, string(privateKey)); err != nil {
		return nil, err
	}
	return plaintext, nil
}

// PriKeyDecrypt 私钥解密
func PriKeyDecrypt(cipherText []byte, privateKey string) ([]byte, error) {
	gRsa := gorsa.RSASecurity{}
	if err := gRsa.SetPrivateKey(privateKey); err != nil {
		return nil, err
	}
	if plainText, err := gRsa.PriKeyDECRYPT(cipherText); err != nil {
		return nil, err
	} else {
		return plainText, nil
	}
}
