package authx

import (
	"github.com/gin-gonic/gin"

	"github.com/go-xuan/quanx/commonx/respx"
	"github.com/go-xuan/quanx/utilx/encryptx"
	"github.com/go-xuan/quanx/utilx/stringx"
)

const (
	AuthType  = "auth-type"     // 鉴权方式
	NoAuth    = "no_auth"       // 免鉴权标识
	Token     = "token"         // token标识
	Cookie    = "cookie"        // cookie鉴权标识
	CookieKey = "auth"          // cookie鉴权key
	TokenKey  = "Authorization" // cookie鉴权key
)

func SetAuthType(ctx *gin.Context, authType string) {
	ctx.Request.Header.Set(AuthType, authType)
}

func SetCookie(ctx *gin.Context, cookie string, age ...int) {
	var maxAge = 3600
	if len(age) > 0 {
		maxAge = age[0]
	}
	if cookie == "" {
		maxAge = -1
	} else {
		var err error
		cookie, err = encryptx.RSA().Encrypt(cookie)
		if err != nil {
			ctx.Abort()
			respx.BuildError(ctx, err)
			return
		}
	}
	ctx.SetCookie(CookieKey, cookie, maxAge, "", "", false, true)
}

// 开启鉴权
func Open() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		switch ctx.Request.Header.Get(AuthType) {
		case NoAuth:
			ctx.Set(AuthType, NoAuth)
		case Token:
			ctx.Set(AuthType, Token)
		default:
			ctx.Set(AuthType, Cookie)
		}
		ctx.Next()
		return
	}
}

// 白名单
func WhiteList(whiteUri ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(whiteUri) > 0 && stringx.ContainsAny(ctx.Request.URL.Path, whiteUri...) {
			ctx.Set(AuthType, NoAuth)
		}
		ctx.Next()
		return
	}
}

// cookie鉴权
func CookeAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetString(AuthType) == Cookie {
			cookie, err := ctx.Cookie(CookieKey)
			if err != nil {
				ctx.Abort()
				respx.Exception(ctx, respx.AuthErr, "cookie is required")
				return
			}
			var account string
			account, err = encryptx.RSA().Decrypt(cookie)
			if err != nil {
				ctx.Abort()
				respx.Exception(ctx, respx.AuthErr, "cookie decrypt failed")
				return
			}
			var user = &User{Account: account}
			if token := user.GetTokenCache(); token == "" {
				ctx.Abort()
				respx.Exception(ctx, respx.AuthErr, "cookie is expired")
				return
			} else {
				user, err = GetUserByToken(token)
				if err != nil {
					ctx.Abort()
					respx.Exception(ctx, respx.AuthErr, "token parse failed")
					return
				}
			}
			ctx.Set("user", user)
			ctx.Next()
		}
	}
}

// token鉴权
func TokenAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetString(AuthType) == Token {
			var token = ctx.Request.Header.Get(TokenKey)
			if token == "" {
				ctx.Abort()
				respx.Exception(ctx, respx.AuthErr, "token is required")
				return
			}
			var user, err = GetUserByToken(token)
			if err != nil || user == nil {
				ctx.Abort()
				respx.Exception(ctx, respx.AuthErr, "token parse failed")
				return
			}
			if user.GetTokenCache() == "" {
				ctx.Abort()
				respx.Exception(ctx, respx.AuthErr, "token is expired")
				return
			}
			ctx.Set("user", user)
			ctx.Next()
		}
	}
}
