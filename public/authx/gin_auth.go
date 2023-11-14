package authx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xuan/quanx/public/redisx"
	"github.com/go-xuan/quanx/public/respx"
	"github.com/go-xuan/quanx/utils/encryptx"
	"github.com/go-xuan/quanx/utils/httpx"
	"github.com/go-xuan/quanx/utils/stringx"
)

const (
	ExemptAuth     = "exempt_auth" // 是否免鉴权
	CookieAuth     = "cookie_auth"
	TokenUser      = "user"
	RedisKeyPrefix = "login@token@"
)

func WhiteList(list ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var username, _, ok = ctx.Request.BasicAuth()
		if username == ExemptAuth && ok {
			ctx.Set(ExemptAuth, true)
		} else {
			ctx.Set(ExemptAuth, stringx.ContainsAny(ctx.Request.RequestURI, list...))
		}
		ctx.Next()
	}
}

// cookie鉴权
func CookeAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetBool(ExemptAuth) {
			return
		}
		cookie, err := ctx.Cookie(CookieAuth)
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
		if token := redisx.GetCmd("user").Get(ctx, RedisKeyPrefix+account).Val(); token == "" {
			ctx.Abort()
			respx.BuildException(ctx, respx.AuthErr, "cookie is expired")
			return
		} else {
			var user *User
			user, err = ParseAuthToken(token)
			if err != nil {
				ctx.Abort()
				respx.BuildException(ctx, respx.AuthErr, "token parse failed")
				return
			}
			ctx.Set(TokenUser, user)
		}
		ctx.Next()
	}
}

// token鉴权
func TokenAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetBool(ExemptAuth) {
			return
		}
		var token = ctx.Request.Header.Get(httpx.Authorization)
		if token == "" {
			ctx.Abort()
			respx.BuildException(ctx, respx.AuthErr, "token is required")
			return
		}
		var user, err = ParseAuthToken(token)
		if err != nil {
			ctx.Abort()
			respx.BuildException(ctx, respx.AuthErr, "token parse failed")
			return
		}
		if redisx.GetCmd("user").Get(ctx, user.RedisCacheKey()).Val() == "" {
			ctx.Abort()
			respx.BuildException(ctx, respx.AuthErr, "token is expired")
			return
		}
		ctx.Set(TokenUser, user)
		ctx.Next()
	}
}
