package ginx

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/go-xuan/quanx/core/cachex"
	"github.com/go-xuan/quanx/net/respx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/encryptx"
)

const (
	AuthType      = "auth_type"     // 鉴权方式
	NoAuth        = "no"            // 免鉴权标识
	Token         = "token"         // token标识
	Cookie        = "cookie"        // cookie鉴权标识
	Authorization = "Authorization" // token鉴权标识
	UserKey       = "user"          // 用户信息存储KEY
	IPKey         = "ip"            // 用户存储KEY
	UserCacheName = "user"          // 用户存储KEY
)

// token 缓存
var cacheClient *cachex.CacheClient

func AuthCache() *cachex.CacheClient {
	if cacheClient == nil {
		cacheClient = cachex.GetClient(UserCacheName)
	}
	return cacheClient
}

func SetAuthType(ctx *gin.Context, authType string) {
	ctx.Request.Header.Set(AuthType, authType)
}

func SetUser(ctx *gin.Context, user *User) {
	ctx.Set(UserKey, user)
}

func GetUser(ctx *gin.Context) *User {
	if value, ok := ctx.Get(UserKey); ok {
		return value.(*User)
	}
	return nil
}

func SetCookie(ctx *gin.Context, username string, age ...int) {
	if cookie, err := encryptx.RSA().Encrypt(username); err != nil {
		ctx.Abort()
		respx.BuildError(ctx, err)
		return
	} else {
		var maxAge = anyx.Default(3600, age...)
		ctx.SetCookie(Cookie, cookie, maxAge, "", "", false, true)
	}
}

// maxAge=-1即可移除cookie
func RemoveCookie(ctx *gin.Context) {
	ctx.SetCookie(Cookie, "", -1, "", "", false, true)
}

// 校验请求IP
func CheckIP(ctx *gin.Context) {
	ctx.Set(IPKey, stringx.IfNot(ctx.ClientIP(), "::1", "localhost"))
	ctx.Next()
	return
}

// 当前请求IP
func GetCorrectIP(ctx *gin.Context) string {
	return ctx.GetString(IPKey)
}

// 开启鉴权
func Auth(ctx *gin.Context) {
	if err := authenticate(ctx); err != nil {
		ctx.Abort()
		respx.Exception(ctx, respx.AuthErr, err.Error())
	} else {
		ctx.Next()
	}
	return
}

func authenticate(ctx *gin.Context) error {
	switch ctx.Request.Header.Get(AuthType) {
	case NoAuth:
		return nil
	case Token:
		return authenticateToken(ctx)
	default:
		return authenticateCookie(ctx)
	}
}

// cookie鉴权
func authenticateToken(ctx *gin.Context) error {
	if token := ctx.Request.Header.Get(Authorization); token == "" {
		return errors.New("token is required")
	} else {
		if user, err := ParseUserFromToken(token); err != nil || user == nil {
			return errors.New("token is invalid")
		} else if !AuthCache().Exists(ctx.Request.Context(), user.Account) {
			return errors.New("token is expired")
		} else {
			SetUser(ctx, user)
		}
		return nil
	}
}

// token鉴权
func authenticateCookie(ctx *gin.Context) error {
	if cookie, err := ctx.Cookie(Cookie); err != nil {
		return errors.New("cookie is required")
	} else {
		var account string
		if account, err = encryptx.RSA().Decrypt(cookie); err == nil {
			if exist := AuthCache().Exists(ctx.Request.Context(), account); !exist {
				return errors.New("cookie is expired")
			} else {
				if token := AuthCache().GetString(ctx.Request.Context(), account); token != "" {
					var user = &User{}
					if user, err = ParseUserFromToken(token); err == nil {
						SetUser(ctx, user)
						return nil
					}
				}
			}
		}
		return errors.New("cookie is invalid")
	}
}
