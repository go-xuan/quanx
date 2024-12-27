package ginx

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/go-xuan/quanx/core/cachex"
	"github.com/go-xuan/quanx/net/respx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/intx"
	"github.com/go-xuan/quanx/utils/encryptx"
)

const (
	tokenHeaderKey = "Authorization"
	sessionUserKey = "SESSION_USER"
	cookieUserName = "COOKIE_USER"
)

var (
	tokenSecret = "123456"         // token加解密密钥，可通过 SetSecret() 方法更改值
	cacheClient cachex.CacheClient // token缓存客户端
)

func AuthValidate() AuthValidator {
	return &JwtValidator{}
}

// AuthValidator 验证器
type AuthValidator interface {
	Token(ctx *gin.Context)  // token鉴权
	Cookie(ctx *gin.Context) // cookie鉴权
	Debug(ctx *gin.Context)  // 本地调试使用
}

// AuthUser 鉴权用户
type AuthUser interface {
	NewToken(secret string) (token string, err error) // 生成token
	ParseToken(token, secret string) error            // 解析token
	UserId() int64                                    // 用户唯一ID
	Username() string                                 // 用户名（唯一）
	Duration() time.Duration                          // 缓存时间
}

func authCache() cachex.CacheClient {
	if cacheClient == nil {
		cacheClient = cachex.Client("auth")
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
	if value, has := ctx.Get(sessionUserKey); has {
		if user, ok := value.(AuthUser); ok {
			return user
		}
	}
	return nil
}

// SetAuthCookie 设置身份验证cookie
func SetAuthCookie(ctx *gin.Context, username string, expire ...int) {
	if cookie, err := encryptx.RSA().Encrypt(username); err != nil {
		ctx.Abort()
		respx.Error(ctx, err.Error())
	} else {
		var maxAge = intx.Default(3600, expire...)
		ctx.SetCookie(cookieUserName, cookie, maxAge, "", "", false, true)
	}
	return
}

// RemoveAuthCookie 移除身份验证cookie maxAge=-1即可移除cookie
func RemoveAuthCookie(ctx *gin.Context) {
	ctx.SetCookie(cookieUserName, "", -1, "", "", false, true)
}

func getSecret() string {
	return tokenSecret
}

func SetSecret(secret string) {
	tokenSecret = secret
}

// SetToken 生成并缓存token
func SetToken(ctx *gin.Context, user AuthUser) (string, error) {
	if token, err := user.NewToken(getSecret()); err != nil {
		return "", errorx.Wrap(err, "new token error")
	} else if err = authCache().Set(ctx, user.Username(), token, user.Duration()); err != nil {
		return "", errorx.Wrap(err, "save token to cache error")
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
	if cookie, err := ctx.Cookie(cookieUserName); err != nil {
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
