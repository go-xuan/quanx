package encryptx

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/go-xuan/quanx/os/filex"
	"os"
	"path"
)

const (
	defaultDir    = "resource/pem"
	privateKeyPem = "rsa-private.pem"
	publicKeyPem  = "rsa-public.pem"
)

var _rsa *Rsa

type Rsa struct {
	privateKey  *rsa.PrivateKey // rsa私钥
	privateData *RsaData        // 私钥数据
	publicData  *RsaData        // 公钥数据
}

type RsaData struct {
	path  string // 秘钥存放路径
	bytes []byte // 秘钥
}

func RSA() *Rsa {
	if _rsa == nil {
		var err error
		if _rsa, err = newRSA(defaultDir); err != nil {
			return nil
		}
	}
	return _rsa
}

// 公钥加密
func (m *Rsa) Encrypt(plaintext string) (ciphertext string, err error) {
	var bytes []byte
	if bytes, err = rsa.EncryptPKCS1v15(rand.Reader, &m.privateKey.PublicKey, []byte(plaintext)); err != nil {
		return
	}
	ciphertext = EncodeBase64(bytes, true)
	return
}

// 私钥解密
func (m *Rsa) Decrypt(ciphertext string) (plaintext string, err error) {
	var base64 = DecodeBase64(ciphertext, true)
	var bytes []byte
	if bytes, err = rsa.DecryptPKCS1v15(rand.Reader, m.privateKey, base64); err != nil {
		return
	}
	plaintext = string(bytes)
	return
}

// 自动生成密钥对并保存到文件
func newRSA(dir string) (r *Rsa, err error) {
	priPath, pubPath := path.Join(dir, privateKeyPem), path.Join(dir, publicKeyPem)
	var priBytes, pubBytes []byte
	var privateKey = &rsa.PrivateKey{}
	if filex.Exists(priPath) && filex.Exists(pubPath) {
		priBytes = PemDecode(priPath)
		if privateKey, err = x509.ParsePKCS1PrivateKey(priBytes); err != nil {
			return
		}
		pubBytes = PemDecode(pubPath)
	} else {
		// 先确保文件夹存在
		filex.CreateDir(dir)
		// 生成私钥
		if privateKey, priBytes, err = NewRSAPrivateKey(priPath); err != nil {
			return
		}
		// 生成公钥
		if pubBytes, err = NewRSAPublicKey(pubPath, &privateKey.PublicKey); err != nil {
			return
		}
	}
	return &Rsa{
		privateKey:  privateKey,
		privateData: &RsaData{priPath, priBytes},
		publicData:  &RsaData{pubPath, pubBytes},
	}, nil
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

// 生成RAS私钥
func NewRSAPrivateKey(pemPath string) (key *rsa.PrivateKey, derBytes []byte, err error) {
	// 生成RSA密钥对
	key = &rsa.PrivateKey{}
	if key, err = rsa.GenerateKey(rand.Reader, 1024); err != nil {
		return
	}
	// 将私钥对象转换成DER编码形式
	derBytes = x509.MarshalPKCS1PrivateKey(key)
	// 创建私钥pem文件
	var file *os.File
	if file, err = os.Create(pemPath); err != nil {
		return
	}
	// 对密钥信息进行编码，写入到私钥文件中
	if err = pem.Encode(file, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: derBytes}); err != nil {
		return
	}
	return
}

// 根据RSA私钥生成公钥
func NewRSAPublicKey(pemPath string, key *rsa.PublicKey) (derBytes []byte, err error) {
	//  将公钥对象序列化为DER编码格式
	derBytes = x509.MarshalPKCS1PublicKey(key)
	// 创建公钥pem文件
	var file *os.File
	if file, err = os.Create(pemPath); err != nil {
		return
	}
	// 对公钥信息进行编码，写入到公钥文件中
	if err = pem.Encode(file, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: derBytes}); err != nil {
		return
	}
	return
}

// RSA加密
func RsaEncrypt(plaintext []byte, pemPath string) (ciphertext []byte, err error) {
	var key *rsa.PublicKey
	if key, err = x509.ParsePKCS1PublicKey(PemDecode(pemPath)); err != nil {
		return
	}
	if ciphertext, err = rsa.EncryptPKCS1v15(rand.Reader, key, plaintext); err != nil {
		return
	}
	return
}

// RSA解密
func RsaDecrypt(ciphertext []byte, pemPath string) (plaintext []byte, err error) {
	var key *rsa.PrivateKey
	if key, err = x509.ParsePKCS1PrivateKey(PemDecode(pemPath)); err != nil {
		return
	}
	if plaintext, err = rsa.DecryptPKCS1v15(rand.Reader, key, ciphertext); err != nil {
		return
	}
	return
}

// RSA加密(PKIX)
func RsaEncryptPKIX(plaintext []byte, pemPath string) (ciphertext []byte, err error) {
	var key any
	if key, err = x509.ParsePKIXPublicKey(PemDecode(pemPath)); err != nil {
		return
	}
	if ciphertext, err = rsa.EncryptPKCS1v15(rand.Reader, key.(*rsa.PublicKey), plaintext); err != nil {
		return
	}
	return
}

// RSA解密(PKIX)
func RsaDecryptPKIX(ciphertext []byte, pemPath string) (plaintext []byte, err error) {
	var key any
	if key, err = x509.ParsePKCS8PrivateKey(PemDecode(pemPath)); err != nil {
		return
	}
	if plaintext, err = rsa.DecryptPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), ciphertext); err != nil {
		return
	}
	return
}
