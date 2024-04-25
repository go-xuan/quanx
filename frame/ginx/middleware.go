package ginx

import (
	"github.com/gin-gonic/gin"

	"github.com/go-xuan/quanx"
	"github.com/go-xuan/quanx/frame/cachex"
	"github.com/go-xuan/quanx/frame/errorx"
	"github.com/go-xuan/quanx/net/respx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/utils/encryptx"
)

const (
	AuthType        = "auth_type"     // 鉴权方式
	NoAuth          = "no"            // 免鉴权标识
	Token           = "token"         // token标识
	Cookie          = "cookie"        // cookie鉴权标识
	Authorization   = "Authorization" // token鉴权标识
	UserKey         = "user"          // 用户信息存储KEY
	IPKey           = "ip"            // 用户存储KEY
	AuthCacheSource = "auth"          // 用户存储KEY
)

// token 缓存
var AuthCache *cachex.CacheClient

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
		var maxAge = anyx.Default(age, 3600)
		ctx.SetCookie(Cookie, cookie, maxAge, "", "", false, true)
	}
}

// maxAge=-1即可移除cookie
func RemoveCookie(ctx *gin.Context) {
	ctx.SetCookie(Cookie, "", -1, "", "", false, true)
}

// 校验请求IP
func CheckIP(ctx *gin.Context) {
	var ip string
	if ip = ctx.ClientIP(); ip == "::1" {
		ip = quanx.GetServer().Host
	}
	ctx.Set(IPKey, ip)
	ctx.Next()
	return
}

// 当前请求IP
func GetCorrectIP(ctx *gin.Context) string {
	return ctx.GetString(IPKey)
}

// 免鉴权
func NotAuth(ctx *gin.Context) {
	SetAuthType(ctx, NoAuth)
	ctx.Next()
	return
}

// 开启鉴权
func Auth(ctx *gin.Context) {
	if AuthCache == nil {
		AuthCache = cachex.Client(AuthCacheSource)
	}
	if err := authenticate(ctx); err != nil {
		ctx.Abort()
		respx.Exception(ctx, respx.AuthErr, err)
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
		return errorx.New("token is required")
	} else {
		if user, err := ParseUserFromToken(token); err != nil || user == nil {
			return errorx.New("token is invalid")
		} else if AuthCache.Get(ctx.Request.Context(), user.Account) == nil {
			return errorx.New("token is expired")
		} else {
			SetUser(ctx, user)
		}
		return nil
	}
}

// token鉴权
func authenticateCookie(ctx *gin.Context) error {
	if cookie, err := ctx.Cookie(Cookie); err != nil {
		return errorx.New("cookie is required")
	} else {
		var account string
		if account, err = encryptx.RSA().Decrypt(cookie); err != nil {
			return errorx.New("cookie is invalid")
		}
		var user = &User{Account: account}
		if token := AuthCache.Get(ctx.Request.Context(), user.Account); token == nil {
			return errorx.New("cookie is expired")
		} else if user, err = ParseUserFromToken(token.(string)); err != nil {
			return errorx.New("cookie is invalid")
		} else {
			SetUser(ctx, user)
		}
		return nil
	}
}
