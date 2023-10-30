package middlewarex

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiaoxuan/quanx/public/redisx"

	"github.com/quanxiaoxuan/quanx/common/constx"
	"github.com/quanxiaoxuan/quanx/common/respx"
	"github.com/quanxiaoxuan/quanx/public/authx"
	"github.com/quanxiaoxuan/quanx/utils/encryptx"
	"github.com/quanxiaoxuan/quanx/utils/stringx"
)

const (
	IsWhiteKey  = "is_white" // 是否在白名单
	CookieKey   = "auth_cookie"
	UserAccount = "user_account"
	UserId      = "user_id"
	UserName    = "user_name"
)

func WhiteList(list ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(list) > 0 {
			ctx.Set(IsWhiteKey, stringx.ContainsAny(ctx.FullPath(), list...))
		}
		ctx.Next()
	}
}

func Cooke() gin.HandlerFunc {
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
		var bytes []byte
		bytes, err = encryptx.RSA().Decrypt([]byte(cookie))
		if err != nil {
			ctx.Abort()
			respx.BuildException(ctx, respx.CookieErr, err.Error())
			return
		}
		var account = string(bytes)
		if redisx.GetCmd("user").Get(ctx, "login:"+account).Val() == "" {
			ctx.Abort()
			respx.BuildException(ctx, respx.AuthErr, "当前用户未登录")
			return
		}
		ctx.Set(UserAccount, account)
		ctx.Next()
	}
}

func Auth() gin.HandlerFunc {
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
		var user, err = authx.ParseAuthToken(token)
		if err != nil {
			ctx.Abort()
			respx.BuildException(ctx, respx.AuthErr, "token解析失败")
			return
		}
		if redisx.GetCmd("user").Get(ctx, "login:"+user.Account).Val() == "" {
			ctx.Abort()
			respx.BuildException(ctx, respx.AuthErr, "当前用户未登录")
			return
		}
		ctx.Set(UserId, user.Id)
		ctx.Set(UserName, user.Name)
		ctx.Next()
	}
}
