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

// RSA 生成RSA加密对象
func RSA(privatePem, publicPem string, mode ...Mode) (*RsaCrypto, error) {
	// 如果私钥和公钥已存在，则直接读取使用，避免重复生成
	if filex.Exists(privatePem) && filex.Exists(publicPem) {
		if rsaCrypto, err := readRsaCrypto(privatePem, publicPem, mode...); err == nil {
			return rsaCrypto, nil
		}
	}
	return newRsaCrypto(privatePem, publicPem, mode...)
}

// RsaCrypto RSA加密器
type RsaCrypto struct {
	privatePem   string          // 私钥文件
	publicPem    string          // 公钥文件
	privateKey   *rsa.PrivateKey // rsa私钥
	privacyBlock *pem.Block      // 私钥块
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
		key, err := x509.ParsePKCS1PrivateKey(c.privacyBlock.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "pkcs1 public key parse error")
		}
		privateKey = key
	case RsaPKCS8:
		key, err := x509.ParsePKCS8PrivateKey(c.privacyBlock.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "pkcs1 public key parse error")
		}
		privateKey = key.(*rsa.PrivateKey)
	default:
		privateKey = c.privateKey
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

// Persistence 密钥持久化
func (c *RsaCrypto) Persistence() error {
	if err := WritePemBlock(c.privatePem, c.privacyBlock); err != nil {
		return errorx.Wrap(err, "save private pem encode error")
	}
	if err := WritePemBlock(c.publicPem, c.publicBlock); err != nil {
		return errorx.Wrap(err, "save public pem error")
	}
	return nil
}

// 生成RAS加密对象
func newRsaCrypto(privatePem, publicPem string, mode ...Mode) (*RsaCrypto, error) {
	var err error

	// 生成RSA密钥
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsa.GenerateKey(rand.Reader, 1024); err != nil {
		return nil, errorx.Wrap(err, "gen rsa private key error")
	}

	// 生成私钥block
	var privateMode, publicMode = parseRsaMode(mode...)
	var privacyBlock, publicBlock *pem.Block
	if privacyBlock, err = GenRSAPrivateBlock(privateKey, privateMode); err != nil {
		return nil, errorx.Wrap(err, "gen rsa private key error")
	}

	// 根据私钥生成公钥block
	if publicBlock, err = GenRSAPublicBlock(privateKey, publicMode); err != nil {
		return nil, errorx.Wrap(err, "gen rsa public key error")
	}

	return &RsaCrypto{
		privatePem:   privatePem,
		publicPem:    publicPem,
		privateKey:   privateKey,
		privacyBlock: privacyBlock,
		publicBlock:  publicBlock,
	}, nil
}

// 从现有的密钥文件中读取RAS加密对象
func readRsaCrypto(privatePem, publicPem string, mode ...Mode) (*RsaCrypto, error) {
	var err error

	// 读取密钥block
	var privacyBlock, publicBlock *pem.Block
	if privacyBlock, err = ReadPemBlock(privatePem); err != nil {
		return nil, errorx.Wrap(err, "read private pem error")
	}
	// 从私钥block解析私钥
	var privateMode, publicMode = parseRsaMode(mode...)
	var privateKey *rsa.PrivateKey
	if privateKey, err = ParseRsaPrivateKey(privacyBlock, privateMode); err != nil {
		return nil, errorx.Wrap(err, "parse private key error")
	}

	// 读取公钥block
	if publicBlock, err = ReadPemBlock(publicPem); err != nil {
		return nil, errorx.Wrap(err, "read public pem error")
	}
	// 从公钥block解析公钥
	if _, err = ParseRsaPublicKey(publicBlock, publicMode); err != nil {
		return nil, errorx.Wrap(err, "parse public key error")
	}

	return &RsaCrypto{
		privatePem:   privatePem,
		publicPem:    publicPem,
		privateKey:   privateKey,
		privacyBlock: privacyBlock,
		publicBlock:  publicBlock,
	}, nil
}

// parseRsaMode 根据传入的模式参数解析RSA私钥和公钥的加解密模式。
func parseRsaMode(mode ...Mode) (Mode, Mode) {
	if len(mode) == 1 && mode[0] == RsaPKCS8 {
		return RsaPKCS8, RsaPKCS1
	} else if len(mode) == 1 && mode[0] == RsaPKIX {
		return RsaPKCS1, RsaPKIX
	} else if len(mode) == 2 {
		return mode[0], mode[1]
	}
	return RsaPKCS1, RsaPKCS1
}

// GenRSAPrivateBlock 生成私钥block
func GenRSAPrivateBlock(privateKey *rsa.PrivateKey, mode Mode) (*pem.Block, error) {
	switch mode {
	case RsaPKCS1:
		data := x509.MarshalPKCS1PrivateKey(privateKey)
		return &pem.Block{
			Type:    "RSA PRIVATE KEY",
			Bytes:   data,
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

// GenRSAPublicBlock 根据私钥生成公钥block
func GenRSAPublicBlock(privateKey *rsa.PrivateKey, mode Mode) (*pem.Block, error) {
	switch mode {
	case RsaPKCS1:
		data := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
		return &pem.Block{
			Type:    "RSA PUBLIC KEY",
			Bytes:   data,
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
		return x509.ParsePKCS1PrivateKey(block.Bytes)
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
		return x509.ParsePKCS1PublicKey(block.Bytes)
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

// WritePemBlock 写入pem文件
func WritePemBlock(path string, block *pem.Block) error {
	var buf bytes.Buffer
	if err := pem.Encode(&buf, block); err != nil {
		return errorx.Wrap(err, "pem encode error")
	}
	if err := filex.WriteFile(path, buf.Bytes()); err != nil {
		return errorx.Wrap(err, "pem write error")
	}
	return nil
}

// ReadPemBlock pem解码
func ReadPemBlock(path string) (*pem.Block, error) {
	data, err := filex.ReadFile(path)
	if err != nil {
		return nil, errorx.Wrap(err, "pem read error")
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errorx.Wrap(err, "pem decode error")
	}
	return block, nil
}
