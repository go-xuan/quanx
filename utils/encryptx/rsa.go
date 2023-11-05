package encryptx

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path"

	"github.com/go-xuan/quanx/utils/filex"
)

const (
	defaultDir    = "pem"
	privateKeyPem = "rsa-private.pem"
	publicKeyPem  = "rsa-public.pem"
)

var myRsa *MyRsa

type MyRsa struct {
	privateKey  *rsa.PrivateKey
	privateData *RsaData
	publicData  *RsaData
}

type RsaData struct {
	path  string
	bytes []byte
}

func RSA() *MyRsa {
	if myRsa == nil {
		var err error
		myRsa, err = newRSA(defaultDir)
		if err != nil {
			return nil
		}
	}
	return myRsa
}

// 公钥加密
func (m *MyRsa) Encrypt(plaintext string) (ciphertext string, err error) {
	var bytes []byte
	bytes, err = rsa.EncryptPKCS1v15(rand.Reader, &m.privateKey.PublicKey, []byte(plaintext))
	if err != nil {
		return
	}
	ciphertext = EncodeBase64(bytes, true)
	return
}

// 私钥解密
func (m *MyRsa) Decrypt(ciphertext string) (plaintext string, err error) {
	var base64 = DecodeBase64(ciphertext, true)
	var bytes []byte
	bytes, err = rsa.DecryptPKCS1v15(rand.Reader, m.privateKey, base64)
	if err != nil {
		return
	}
	plaintext = string(bytes)
	return
}

// 生成密钥对并保存到文件
func newRSA(dir string) (myRsa *MyRsa, err error) {
	priPath, pubPath := path.Join(dir, privateKeyPem), path.Join(dir, publicKeyPem)
	var priBytes, pubBytes []byte
	var privateKey = &rsa.PrivateKey{}
	if filex.Exists(priPath) && filex.Exists(pubPath) {
		priBytes = PemDecode(priPath)
		privateKey, err = x509.ParsePKCS1PrivateKey(priBytes)
		if err != nil {
			return
		}
		pubBytes = PemDecode(pubPath)
	} else {
		// 先确保文件夹存在
		filex.CreateDir(dir)
		// 生成私钥
		privateKey, priBytes, err = NewRSAPrivateKey(priPath)
		if err != nil {
			return
		}
		// 生成公钥
		pubBytes, err = NewRSAPublicKey(pubPath, &privateKey.PublicKey)
		if err != nil {
			return
		}
	}
	myRsa = &MyRsa{
		privateKey:  privateKey,
		privateData: &RsaData{priPath, priBytes},
		publicData:  &RsaData{pubPath, pubBytes},
	}
	return
}

// pem解码
func PemDecode(pemPath string) []byte {
	if data, err := os.ReadFile(pemPath); err == nil {
		if block, _ := pem.Decode(data); block != nil {
			return block.Bytes
		}
	}
	return nil
}

func NewRSAPrivateKey(pemPath string) (key *rsa.PrivateKey, derBytes []byte, err error) {
	// 生成RSA密钥对
	key = &rsa.PrivateKey{}
	key, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return
	}
	// 将私钥对象转换成DER编码形式
	derBytes = x509.MarshalPKCS1PrivateKey(key)
	// 创建私钥pem文件
	var file *os.File
	file, err = os.Create(pemPath)
	if err != nil {
		return
	}
	// 对密钥信息进行编码，写入到私钥文件中
	err = pem.Encode(file, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: derBytes})
	if err != nil {
		return
	}
	return
}

func NewRSAPublicKey(pemPath string, key *rsa.PublicKey) (derBytes []byte, err error) {
	//  将公钥对象序列化为DER编码格式
	derBytes = x509.MarshalPKCS1PublicKey(key)
	// 创建公钥pem文件
	var file *os.File
	file, err = os.Create(pemPath)
	if err != nil {
		return
	}
	// 对公钥信息进行编码，写入到公钥文件中
	err = pem.Encode(file, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: derBytes})
	if err != nil {
		return
	}
	return
}

// 根据第三方公钥加密
func RsaEncrypt(plaintext []byte, pemPath string) (ciphertext []byte, err error) {
	var key *rsa.PublicKey
	key, err = x509.ParsePKCS1PublicKey(PemDecode(pemPath))
	if err != nil {
		return
	}
	ciphertext, err = rsa.EncryptPKCS1v15(rand.Reader, key, plaintext)
	if err != nil {
		return
	}
	return
}

// 根据第三方私钥解密
func RsaDecrypt(plaintext []byte, pemPath string) (ciphertext []byte, err error) {
	var key *rsa.PrivateKey
	key, err = x509.ParsePKCS1PrivateKey(PemDecode(pemPath))
	if err != nil {
		return
	}
	ciphertext, err = rsa.DecryptPKCS1v15(rand.Reader, key, plaintext)
	if err != nil {
		return
	}
	return
}

// 根据第三方公钥加密
func RsaEncryptPKIX(plaintext []byte, pemPath string) (ciphertext []byte, err error) {
	var key any
	key, err = x509.ParsePKIXPublicKey(PemDecode(pemPath))
	if err != nil {
		return
	}
	ciphertext, err = rsa.EncryptPKCS1v15(rand.Reader, key.(*rsa.PublicKey), plaintext)
	if err != nil {
		return
	}
	return
}

// 根据第三方私钥解密
func RsaDecryptPKIX(plaintext []byte, pemPath string) (ciphertext []byte, err error) {
	var key any
	key, err = x509.ParsePKCS8PrivateKey(PemDecode(pemPath))
	if err != nil {
		return
	}
	ciphertext, err = rsa.DecryptPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), plaintext)
	if err != nil {
		return
	}
	return
}
