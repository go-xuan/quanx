package authx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xuan/quanx/utils/defaultx"
	"github.com/go-xuan/quanx/utils/encryptx"

	"github.com/go-xuan/quanx/common/respx"
	"github.com/go-xuan/quanx/utils/stringx"
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

func WhiteList(uri ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqUrl = ctx.Request.URL.Path
		var authType = ctx.Request.Header.Get(AuthType)
		if authType == NoAuth {
			ctx.Set(AuthType, NoAuth)
		} else if stringx.ContainsAny(reqUrl, uri...) {
			ctx.Set(AuthType, NoAuth)
		} else {
			ctx.Set(AuthType, defaultx.String(authType, Cookie))
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
				respx.BuildException(ctx, respx.AuthErr, "cookie is required")
				return
			}
			var account string
			account, err = encryptx.RSA().Decrypt(cookie)
			if err != nil {
				ctx.Abort()
				respx.BuildException(ctx, respx.AuthErr, "cookie decrypt failed")
				return
			}
			var user = &User{Account: account}
			if token := user.GetTokenCache(); token == "" {
				ctx.Abort()
				respx.BuildException(ctx, respx.AuthErr, "cookie is expired")
				return
			} else {
				user, err = GetUserByToken(token)
				if err != nil {
					ctx.Abort()
					respx.BuildException(ctx, respx.AuthErr, "token parse failed")
					return
				}
				ctx.Set("user", user)
			}
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
				respx.BuildException(ctx, respx.AuthErr, "token is required")
				return
			}
			var user, err = GetUserByToken(token)
			if err != nil || user == nil {
				ctx.Abort()
				respx.BuildException(ctx, respx.AuthErr, "token parse failed")
				return
			}
			if user.GetTokenCache() == "" {
				ctx.Abort()
				respx.BuildException(ctx, respx.AuthErr, "token is expired")
				return
			}
			ctx.Set("user", user)
			ctx.Next()
		}
	}
}
