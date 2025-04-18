package cryptx

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

// MD5 MD5加密
func MD5(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

// SHA1 sha1加密
func SHA1(s string) string {
	hash := sha1.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

// PasswordSalt 密码加盐
func PasswordSalt(password, salt string) string {
	hash := hmac.New(sha1.New, []byte(salt))
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// Base64Encode base64加密
func Base64Encode(text []byte, safe ...bool) string {
	if len(safe) > 0 && safe[0] {
		return base64.URLEncoding.EncodeToString(text)
	} else {
		return base64.StdEncoding.EncodeToString(text)
	}
}

// Base64Decode base64解密
func Base64Decode(text string, safe ...bool) ([]byte, error) {
	if len(safe) > 0 && safe[0] {
		return base64.URLEncoding.DecodeString(text)
	} else {
		return base64.StdEncoding.DecodeString(text)
	}
}
