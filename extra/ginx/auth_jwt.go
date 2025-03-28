package ginx

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/respx"
	"github.com/go-xuan/quanx/types/anyx"
)

type JwtValidator struct{}

func (v JwtValidator) Token(ctx *gin.Context) {
	if err := ValidateToken(ctx, &JwtUser{}); err != nil {
		respx.Forbidden(ctx, err)
		ctx.Abort()
	}
}

func (v JwtValidator) Cookie(ctx *gin.Context) {
	if err := ValidateCookie(ctx, &JwtUser{}); err != nil {
		respx.Forbidden(ctx, err)
		ctx.Abort()
	}
}

func (v JwtValidator) Debug(ctx *gin.Context) {
	SetSessionUser(ctx, &JwtUser{Id: 999999999})
}

var jwtSecret = "123456" // token加解密密钥，可通过 SetSecret() 方法更改值

func getJwtSecret() interface{} {
	return []byte(jwtSecret)
}

func jwtKeyFunc(*jwt.Token) (interface{}, error) {
	return getJwtSecret(), nil
}

// JwtUser jwt-TokenUser实现
type JwtUser struct {
	Id     int64 `json:"id"`     // 用户ID
	Expire int64 `json:"expire"` // 存活时长
}

func (u *JwtUser) Valid() error {
	if u.Expire < time.Now().Unix() {
		return errorx.New("claims valid failed: current user expired")
	}
	return nil
}

func (u *JwtUser) Encrypt() (string, error) {
	if token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, u).SignedString(getJwtSecret()); err != nil {
		return "", errorx.Wrap(err, "jwt encrypt error")
	} else {
		return token, nil
	}
}

func (u *JwtUser) Decrypt(token string) error {
	if jt, err := jwt.ParseWithClaims(token, u, jwtKeyFunc); err != nil || !jt.Valid {
		return errorx.Wrap(err, "jwt decrypt error")
	}
	return nil
}

func (u *JwtUser) UserId() anyx.Value {
	return anyx.Int64Value(u.Id)
}

func (u *JwtUser) TTL() time.Duration {
	return 3600 * time.Second
}
