package ginx

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/go-xuan/quanx/core/cachex"
	"github.com/go-xuan/quanx/net/respx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/intx"
	"github.com/go-xuan/quanx/utils/encryptx"
)

const (
	tokenHeaderKey = "Authorization"
	sessionUserKey = "session_user"
	cookieName     = "cookie_user"
)

func AuthValidate() AuthValidator {
	return &JwtValidator{}
}

// AuthValidator 验证器接口
type AuthValidator interface {
	Token(ctx *gin.Context)  // token鉴权
	Cookie(ctx *gin.Context) // cookie鉴权
	Debug(ctx *gin.Context)  // 本地调试使用
}

// AuthUser 鉴权用户接口
type AuthUser interface {
	NewToken(secret string) (token string, err error) // 生成token
	ParseToken(token, secret string) error            // 解析token
	UserId() int64                                    // 用户唯一ID
	Username() string                                 // 用户名（唯一）
	Duration() time.Duration                          // 缓存时间
}

// token 缓存
var cacheClient cachex.Client

func authCache() cachex.Client {
	if cacheClient == nil {
		cacheClient = cachex.GetClient("auth")
	}
	return cacheClient
}

// SetSessionUser 设置会话用户
func SetSessionUser(ctx *gin.Context, user AuthUser) {
	ctx.Set(sessionUserKey, user)
	// token续命
	_ = authCache().Expire(ctx, user.Username(), user.Duration())
}

// GetSessionUser 获取会话用户
func GetSessionUser(ctx *gin.Context) AuthUser {
	if value, ok := ctx.Get(sessionUserKey); ok {
		return value.(AuthUser)
	}
	return nil
}

// SetAuthCookie 设置身份验证cookie
func SetAuthCookie(ctx *gin.Context, username string, expire ...int) {
	if cookie, err := encryptx.RSA().Encrypt(username); err != nil {
		ctx.Abort()
		respx.Ctx(ctx).Error(err)
	} else {
		var maxAge = intx.Default(3600, expire...)
		ctx.SetCookie(cookieName, cookie, maxAge, "", "", false, true)
	}
	return
}

// RemoveAuthCookie 移除身份验证cookie maxAge=-1即可移除cookie
func RemoveAuthCookie(ctx *gin.Context) {
	ctx.SetCookie(cookieName, "", -1, "", "", false, true)
}

// token加解密密钥
var tokenSecret = "123456"

func getSecret() string {
	return tokenSecret
}

func SetSecret(secret string) {
	tokenSecret = secret
}

// SetToken 生成并缓存token
func SetToken(ctx *gin.Context, user AuthUser) (string, error) {
	if token, err := user.NewToken(getSecret()); err != nil {
		return "", errorx.Wrap(err, "生成token失败")
	} else if err = authCache().Set(ctx, user.Username(), token, user.Duration()); err != nil {
		return "", errorx.Wrap(err, "token写入缓存失败")
	} else {
		return token, nil
	}
}

// RemoveToken 移除token
func RemoveToken(ctx *gin.Context, username string) {
	authCache().Delete(ctx, username)
}

// token方式鉴权
func authValidateWithToken(ctx *gin.Context, user AuthUser) error {
	if token := ctx.Request.Header.Get(tokenHeaderKey); token == "" {
		return errorx.New("token is required")
	} else {
		if err := user.ParseToken(token, getSecret()); err != nil {
			return errorx.Wrap(err, "parse token failed")
		}
		if exist := authCache().Exist(ctx, user.Username()); !exist {
			return errorx.New("token has expired")
		}
		SetSessionUser(ctx, user)
		return nil
	}
}

// cookie方式鉴权
func authValidateWithCookie(ctx *gin.Context, user AuthUser) error {
	if cookie, err := ctx.Cookie(cookieName); err != nil {
		return errorx.Wrap(err, "get request cookie failed")
	} else {
		var username string
		if username, err = encryptx.RSA().Decrypt(cookie); err != nil {
			return errorx.Wrap(err, "cookie is invalid")
		}
		var token string
		if exist := authCache().Get(ctx, username, &token); !exist || token == "" {
			return errorx.Errorf("cookie has expired: %s", username)
		}
		if err = user.ParseToken(token, getSecret()); err != nil {
			return errorx.Wrap(err, "parse token failed")
		}
		SetSessionUser(ctx, user)
		return nil
	}
}

// JwtValidator JWT鉴权验证器
type JwtValidator struct{}

func (v *JwtValidator) Token(ctx *gin.Context) {
	if err := authValidateWithToken(ctx, &JwtUser{}); err != nil {
		ctx.Abort()
		respx.Ctx(ctx).Forbidden(err)
	} else {
		ctx.Next()
	}
	return
}

func (v *JwtValidator) Cookie(ctx *gin.Context) {
	if err := authValidateWithCookie(ctx, &JwtUser{}); err != nil {
		ctx.Abort()
		respx.Ctx(ctx).Forbidden(err)
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

func (u *JwtUser) NewToken(secret string) (token string, err error) {
	if token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, u).SignedString([]byte(secret)); err != nil {
		return
	}
	return
}

func (u *JwtUser) ParseToken(token, secret string) (err error) {
	var jwtToken *jwt.Token
	if jwtToken, err = jwt.ParseWithClaims(token, &JwtUser{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}); err != nil {
		err = errorx.Wrap(err, "解析token失败")
		return
	}
	if user, ok := jwtToken.Claims.(*JwtUser); ok {
		u.Id = user.Id
		u.Account = user.Account
		u.Name = user.Name
		u.Phone = user.Phone
		u.Ip = user.Ip
		u.Domain = user.Domain
		u.TTL = user.TTL
	}
	return
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
