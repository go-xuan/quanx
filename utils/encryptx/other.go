package encryptx

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

// EncodeBase64 base64加密
func EncodeBase64(text []byte, safe ...bool) (result string) {
	if len(safe) > 0 && safe[0] {
		result = base64.URLEncoding.EncodeToString(text)
	} else {
		result = base64.StdEncoding.EncodeToString(text)
	}
	return
}

// DecodeBase64 base64解密
func DecodeBase64(text string, safe ...bool) (result []byte) {
	if len(safe) > 0 && safe[0] {
		result, _ = base64.URLEncoding.DecodeString(text)
	} else {
		result, _ = base64.StdEncoding.DecodeString(text)
	}
	return
}
