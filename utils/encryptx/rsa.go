package encryptx

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
)

const (
	defaultDir    = "rsa"
	privateKeyPem = "rsa-private.pem"
	publicKeyPem  = "rsa-public.pem"
)

var _rsa *Rsa

type Rsa struct {
	privateKey  *rsa.PrivateKey // rsa私钥
	privatePath string          // 秘钥存放路径
	privateData []byte          // 秘钥
	publicPath  string          // 公钥存放路径
	publicData  []byte          // 公钥
}

type RsaData struct {
	path  string // 秘钥存放路径
	bytes []byte // 秘钥
}

func RSA() *Rsa {
	if _rsa == nil {
		var err error
		if _rsa, err = newRSA(defaultDir); err != nil {
			panic(err)
		}
	}
	return _rsa
}

// Encrypt 公钥加密
func (m *Rsa) Encrypt(plaintext string) (string, error) {
	if ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, &m.privateKey.PublicKey, []byte(plaintext)); err != nil {
		return "", errorx.Wrap(err, "pkcs1v15 encrypt error")
	} else {
		return Base64Encode(ciphertext, true), nil
	}
}

// Decrypt 私钥解密
func (m *Rsa) Decrypt(ciphertext string) (string, error) {
	if base64, err := Base64Decode(ciphertext, true); err != nil {
		return "", errorx.Wrap(err, "base64 decode error")
	} else {
		var plaintext []byte
		if plaintext, err = rsa.DecryptPKCS1v15(rand.Reader, m.privateKey, base64); err != nil {
			return "", errorx.Wrap(err, "pkcs1v15 decrypt error")
		}
		return string(plaintext), nil
	}
}

// 自动生成密钥对并保存到文件
func newRSA(dir string) (*Rsa, error) {
	priPath, pubPath := filepath.Join(dir, privateKeyPem), filepath.Join(dir, publicKeyPem)
	if filex.Exists(priPath) && filex.Exists(pubPath) {
		var priBytes = pemDecode(priPath)
		if privateKey, err := x509.ParsePKCS1PrivateKey(priBytes); err != nil {
			return nil, errorx.Wrap(err, "pkcs1v15 private key parse error")
		} else {
			return &Rsa{
				privateKey:  privateKey,
				privatePath: priPath,
				privateData: priBytes,
				publicPath:  pubPath,
				publicData:  pemDecode(pubPath),
			}, nil
		}
	} else {
		// 先确保文件夹存在
		filex.CreateDir(dir)
		// 生成私钥
		if privateKey, priBytes, err := NewRSAPrivateKey(priPath); err != nil {
			return nil, errorx.Wrap(err, "new rsa private key error")
		} else {
			// 根据私钥生成公钥
			var pubBytes []byte
			if pubBytes, err = NewRSAPublicKey(pubPath, &privateKey.PublicKey); err != nil {
				return nil, errorx.Wrap(err, "new rsa public key error")
			}
			return &Rsa{
				privateKey:  privateKey,
				privatePath: priPath,
				privateData: priBytes,
				publicPath:  pubPath,
				publicData:  pubBytes,
			}, nil
		}
	}
}

// pemDecode pem解码
func pemDecode(path string) []byte {
	if data, err := os.ReadFile(path); err == nil {
		if block, _ := pem.Decode(data); block != nil {
			return block.Bytes
		}
	}
	return nil
}

// NewRSAPrivateKey 生成RAS私钥
func NewRSAPrivateKey(path string) (*rsa.PrivateKey, []byte, error) {
	// 生成RSA密钥对
	if key, err := rsa.GenerateKey(rand.Reader, 1024); err != nil {
		return nil, nil, err
	} else {
		// 创建私钥pem文件
		var file *os.File
		if file, err = os.Create(path); err != nil {
			return nil, nil, err
		}
		defer file.Close()
		// 将私钥对象转换成DER编码形式
		derBytes := x509.MarshalPKCS1PrivateKey(key)
		// 对密钥信息进行编码，写入到私钥文件中
		if err = pem.Encode(file, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: derBytes}); err != nil {
			return nil, nil, err
		}
		return key, derBytes, nil
	}
}

// NewRSAPublicKey 根据RSA私钥生成公钥
func NewRSAPublicKey(path string, key *rsa.PublicKey) ([]byte, error) {
	// 创建公钥pem文件
	if pemFile, err := os.Create(path); err != nil {
		return nil, err
	} else {
		defer pemFile.Close()
		//  将公钥对象序列化为DER编码格式
		var derBytes = x509.MarshalPKCS1PublicKey(key)
		// 对公钥信息进行编码，写入到公钥文件中
		if err = pem.Encode(pemFile, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: derBytes}); err != nil {
			return nil, err
		}
		return derBytes, nil
	}
}

// RsaEncrypt RSA加密
func RsaEncrypt(plaintext []byte, pubPem string) ([]byte, error) {
	if publicKey, err := x509.ParsePKCS1PublicKey(pemDecode(pubPem)); err != nil {
		return nil, errorx.Wrap(err, "pkcs1 public key parse error")
	} else {
		return rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
	}
}

// RsaDecrypt RSA解密
func RsaDecrypt(ciphertext []byte, priPem string) ([]byte, error) {
	if privateKey, err := x509.ParsePKCS1PrivateKey(pemDecode(priPem)); err != nil {
		return nil, err
	} else {
		return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	}
}

// RsaEncryptPKIX RSA加密(PKIX)
func RsaEncryptPKIX(plaintext []byte, pubPem string) ([]byte, error) {
	if publicKey, err := x509.ParsePKIXPublicKey(pemDecode(pubPem)); err != nil {
		return nil, err
	} else {
		return rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), plaintext)
	}
}

// RsaDecryptPKIX RSA解密(PKIX)
func RsaDecryptPKIX(ciphertext []byte, priPem string) ([]byte, error) {
	if privateKey, err := x509.ParsePKCS8PrivateKey(pemDecode(priPem)); err != nil {
		return nil, err
	} else {
		return rsa.DecryptPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), ciphertext)
	}
}
