package cryptox

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strconv"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
)

const PemHeaderModeKey = "Crypto-Mode"

func RSA() (*RsaCrypto, error) {
	var crypto *RsaCrypto
	if privateKey, err := rsa.GenerateKey(rand.Reader, 1024); err != nil {
		return nil, err
	} else if crypto, err = NewRsaCrypto(privateKey, RsaPKCS1, RsaPKCS1); err != nil {
		return nil, err
	}
	return crypto, nil
}

// NewRsaCrypto 生成RAS加密对象
func NewRsaCrypto(privateKey *rsa.PrivateKey, privateMode, publicMode Mode) (*RsaCrypto, error) {
	var err error

	// 生成私钥block
	var privacyBlock *pem.Block
	if privacyBlock, err = NewRSAPrivateBlock(privateKey, privateMode); err != nil {
		return nil, errorx.Wrap(err, "gen rsa private key error")
	}
	// 根据私钥生成公钥block
	var publicBlock *pem.Block
	if publicBlock, err = NewRSAPublicBlock(privateKey, publicMode); err != nil {
		return nil, errorx.Wrap(err, "gen rsa public key error")
	}
	return &RsaCrypto{
		privateKey:   privateKey,
		privateBlock: privacyBlock,
		publicBlock:  publicBlock,
	}, nil
}

// ParseRsaCrypto 从现有的密钥解析出RAS加密对象
func ParseRsaCrypto(privateData []byte, privateMode, publicMode Mode) (*RsaCrypto, error) {
	var err error

	var privateBlock *pem.Block
	if privateBlock, _ = pem.Decode(privateData); privateBlock == nil {
		return nil, errorx.Wrap(err, "private key decode error")
	}
	var privateKey *rsa.PrivateKey
	if privateKey, err = ParseRsaPrivateKey(privateBlock, privateMode); err != nil {
		return nil, errorx.Wrap(err, "parse private key error")
	}
	var publicBlock *pem.Block
	if publicBlock, err = NewRSAPublicBlock(privateKey, publicMode); err != nil {
		return nil, errorx.Wrap(err, "gen rsa public key error")
	}
	return &RsaCrypto{
		privateKey:   privateKey,
		privateBlock: privateBlock,
		publicBlock:  publicBlock,
	}, nil
}

// RsaCrypto RSA加密器
type RsaCrypto struct {
	privateKey   *rsa.PrivateKey // rsa私钥
	privateBlock *pem.Block      // 私钥块
	publicBlock  *pem.Block      // 公钥块
}

// Encrypt 公钥加密
func (c *RsaCrypto) Encrypt(plaintext []byte) ([]byte, error) {
	var publicKey *rsa.PublicKey
	mode, _ := strconv.Atoi(c.publicBlock.Headers[PemHeaderModeKey])
	switch Mode(mode) {
	case RsaPKCS1:
		key, err := x509.ParsePKCS1PublicKey(c.publicBlock.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "pkcs1 public key parse error")
		}
		publicKey = key
	case RsaPKIX:
		key, err := x509.ParsePKIXPublicKey(c.publicBlock.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "pkcs1 public key parse error")
		}
		publicKey = key.(*rsa.PublicKey)
	default:
		publicKey = &c.privateKey.PublicKey
	}
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
}

// Decrypt 私钥解密
func (c *RsaCrypto) Decrypt(ciphertext []byte) ([]byte, error) {
	var privateKey *rsa.PrivateKey
	mode, _ := strconv.Atoi(c.publicBlock.Headers[PemHeaderModeKey])
	switch Mode(mode) {
	case RsaPKCS1:
		key, err := x509.ParsePKCS1PrivateKey(c.privateBlock.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "pkcs1 public key parse error")
		}
		privateKey = key
	case RsaPKCS8:
		key, err := x509.ParsePKCS8PrivateKey(c.privateBlock.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "pkcs1 public key parse error")
		}
		privateKey = key.(*rsa.PrivateKey)
	default:
		privateKey = c.privateKey
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

// SavePrivateKey 密钥持久化
func (c *RsaCrypto) SavePrivateKey(path string) error {
	if err := writePemBlock(path, c.privateBlock); err != nil {
		return errorx.Wrap(err, "save private pem encode error")
	}
	return nil
}

// SavePublicKey 公钥持久化
func (c *RsaCrypto) SavePublicKey(path string) error {
	if err := writePemBlock(path, c.publicBlock); err != nil {
		return errorx.Wrap(err, "save public pem error")
	}
	return nil
}

// NewRSAPrivateBlock 生成私钥block
func NewRSAPrivateBlock(privateKey *rsa.PrivateKey, mode Mode) (*pem.Block, error) {
	switch mode {
	case RsaPKCS1:
		return &pem.Block{
			Type:    "RSA PRIVATE KEY",
			Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
			Headers: map[string]string{PemHeaderModeKey: strconv.Itoa(int(RsaPKCS1))},
		}, nil
	case RsaPKCS8:
		data, err := x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			return nil, errorx.Wrap(err, "marshal PKCS8 private key error")
		}
		return &pem.Block{
			Type:    "RSA PRIVATE KEY",
			Bytes:   data,
			Headers: map[string]string{PemHeaderModeKey: strconv.Itoa(int(RsaPKCS8))},
		}, nil
	default:
		return nil, errorx.New("unsupported private mode:" + strconv.Itoa(int(mode)))
	}
}

// NewRSAPublicBlock 根据私钥生成公钥block
func NewRSAPublicBlock(privateKey *rsa.PrivateKey, mode Mode) (*pem.Block, error) {
	switch mode {
	case RsaPKCS1:
		return &pem.Block{
			Type:    "RSA PUBLIC KEY",
			Bytes:   x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
			Headers: map[string]string{"mode": strconv.Itoa(int(RsaPKCS1))},
		}, nil
	case RsaPKIX:
		data, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		if err != nil {
			return nil, errorx.Wrap(err, "marshal PKIX private privateKey error")
		}
		return &pem.Block{
			Type:    "RSA PUBLIC KEY",
			Bytes:   data,
			Headers: map[string]string{"mode": strconv.Itoa(int(RsaPKIX))},
		}, nil
	default:
		return nil, errorx.New("unsupported public mode:" + strconv.Itoa(int(mode)))
	}
}

// ParseRsaPrivateKey 从私钥block解析私钥
func ParseRsaPrivateKey(block *pem.Block, mode Mode) (*rsa.PrivateKey, error) {
	switch mode {
	case RsaPKCS1:
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "parse PKCS1 private key error")
		}
		return key, nil
	case RsaPKCS8:
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "parse PKCS8 private key error")
		}
		return key.(*rsa.PrivateKey), nil
	default:
		return nil, errorx.New("unsupported private mode:" + strconv.Itoa(int(mode)))
	}
}

// ParseRsaPublicKey 从公钥block解析公钥
func ParseRsaPublicKey(block *pem.Block, mode Mode) (*rsa.PublicKey, error) {
	switch mode {
	case RsaPKCS1:
		key, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "parse PKCS1 public key error")
		}
		return key, nil
	case RsaPKIX:
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "parse PKIX private key error")
		}
		return key.(*rsa.PublicKey), nil
	default:
		return nil, errorx.New("unsupported public mode:" + strconv.Itoa(int(mode)))
	}
}

// writePemBlock 写入pem文件
func writePemBlock(path string, block *pem.Block) error {
	var buf bytes.Buffer
	if err := pem.Encode(&buf, block); err != nil {
		return errorx.Wrap(err, "pem encode error")
	}
	if err := filex.WriteFile(path, buf.Bytes()); err != nil {
		return errorx.Wrap(err, "pem write error")
	}
	return nil
}
