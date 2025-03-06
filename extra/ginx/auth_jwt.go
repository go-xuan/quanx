package ginx

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/respx"
	"github.com/go-xuan/quanx/types/intx"
)

// JwtValidator JWT鉴权验证器
type JwtValidator struct{}

func (v *JwtValidator) Token(ctx *gin.Context) {
	if err := authValidateWithToken(ctx, &JwtUser{}); err != nil {
		respx.Forbidden(ctx, err)
		ctx.Abort()
	} else {
		ctx.Next()
	}
	return
}

func (v *JwtValidator) Cookie(ctx *gin.Context) {
	if err := authValidateWithCookie(ctx, &JwtUser{}); err != nil {
		respx.Forbidden(ctx, err)
		ctx.Abort()
	} else {
		ctx.Next()
	}
	return
}

func (v *JwtValidator) Debug(ctx *gin.Context) {
	SetSessionUser(ctx, &JwtUser{
		Id:      999999999,
		Account: "debug",
		Name:    "debug",
		Phone:   "110",
	})
	ctx.Next()
	return
}

// JwtUser jwt-TokenUser实现
type JwtUser struct {
	Id      int64  `json:"id"`      // 用户ID
	Account string `json:"account"` // 用户账号
	Name    string `json:"name"`    // 用户姓名
	Phone   string `json:"phone"`   // 登录手机
	Ip      string `json:"ip"`      // 登录IP
	Domain  string `json:"domain"`  // 域名
	TTL     int    `json:"ttl"`     // 有效时长
}

func (u *JwtUser) Valid() error {
	return nil
}

func (u *JwtUser) NewToken(secret string) (string, error) {
	if token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, u).SignedString([]byte(secret)); err != nil {
		return "", errorx.Wrap(err, "token sign failed")
	} else {
		return token, nil
	}
}

func (u *JwtUser) ParseToken(token, secret string) error {
	if jwtToken, err := jwt.ParseWithClaims(token, &JwtUser{}, func(*jwt.Token) (any, error) {
		return []byte(secret), nil
	}); err != nil {
		return errorx.Wrap(err, "parse token error")
	} else if user, ok := jwtToken.Claims.(*JwtUser); ok {
		u.Id = user.Id
		u.Account = user.Account
		u.Name = user.Name
		u.Phone = user.Phone
		u.Ip = user.Ip
		u.Domain = user.Domain
		u.TTL = user.TTL
	}
	return nil
}

func (u *JwtUser) Username() string {
	return u.Phone
}

func (u *JwtUser) UserId() int64 {
	return u.Id
}

func (u *JwtUser) Duration() time.Duration {
	return time.Duration(intx.IfZero(u.TTL, 3600)) * time.Second
}
