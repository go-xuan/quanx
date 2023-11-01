package authx

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiaoxuan/quanx/common/constx"
	"github.com/quanxiaoxuan/quanx/common/respx"
	"github.com/quanxiaoxuan/quanx/public/redisx"
	"github.com/quanxiaoxuan/quanx/utils/encryptx"
	"github.com/quanxiaoxuan/quanx/utils/stringx"
)

const (
	IsWhiteKey     = "is_white" // 是否在白名单
	CookieKey      = "cookie_auth"
	TokenUser      = "user"
	RedisKeyPrefix = "login@token@"
)

func WhiteList(list ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(list) > 0 {
			ctx.Set(IsWhiteKey, stringx.ContainsAny(ctx.FullPath(), list...))
		}
		ctx.Next()
	}
}

// cookie鉴权
func CookeAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetBool(IsWhiteKey) {
			return
		}
		cookie, err := ctx.Cookie(CookieKey)
		if err != nil {
			ctx.Abort()
			respx.BuildException(ctx, respx.CookieErr, err.Error())
			return
		}
		var account string
		account, err = encryptx.RSA().Decrypt(cookie)
		if err != nil {
			ctx.Abort()
			respx.BuildException(ctx, respx.CookieErr, err.Error())
			return
		}
		if token := redisx.GetCmd("user").Get(ctx, RedisKeyPrefix+account).Val(); token == "" {
			ctx.Abort()
			respx.BuildException(ctx, respx.AuthErr, "当前cookie已失效")
			return
		} else {
			var user *User
			user, err = ParseAuthToken(token)
			if err != nil {
				ctx.Abort()
				respx.BuildException(ctx, respx.AuthErr, "token解析失败")
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
		if ctx.GetBool(IsWhiteKey) {
			return
		}
		var token = ctx.Request.Header.Get(constx.Authorization)
		if token == "" {
			ctx.Abort()
			respx.BuildException(ctx, respx.AuthErr, "未携带token")
			return
		}
		var user, err = ParseAuthToken(token)
		if err != nil {
			ctx.Abort()
			respx.BuildException(ctx, respx.AuthErr, "token解析失败")
			return
		}
		if redisx.GetCmd("user").Get(ctx, user.RedisCacheKey()).Val() == "" {
			ctx.Abort()
			respx.BuildException(ctx, respx.AuthErr, "当前token已失效")
			return
		}
		ctx.Set(TokenUser, user)
		ctx.Next()
	}
}
