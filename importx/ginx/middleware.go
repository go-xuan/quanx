package ginx

import (
	"github.com/gin-gonic/gin"

	"github.com/go-xuan/quanx/commonx/respx"
	"github.com/go-xuan/quanx/utilx/encryptx"
	"github.com/go-xuan/quanx/utilx/ipx"
	"github.com/go-xuan/quanx/utilx/stringx"
)

const (
	AuthType = "auth-type" // 鉴权方式
	NoAuth   = "no_auth"   // 免鉴权标识
	Token    = "token"     // token标识
	Cookie   = "cookie"    // cookie鉴权标识
)

func RemoveCookie(ctx *gin.Context) {
	ctx.SetCookie(Cookie, "", -1, "", "", false, true)
}

func SetCookie(ctx *gin.Context, username string, age ...int) {
	var maxAge = 3600
	if len(age) > 0 {
		maxAge = age[0]
	}
	cookie, err := encryptx.RSA().Encrypt(username)
	if err != nil {
		ctx.Abort()
		respx.BuildError(ctx, err)
		return
	}
	ctx.SetCookie(Cookie, cookie, maxAge, "", "", false, true)
}

func IP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		if ip == "::1" {
			ip = ipx.GetWLANIP()
		}
		ctx.Set("ip", ip)
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
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var msg string
		switch ctx.Request.Header.Get(AuthType) {
		case NoAuth:
			msg = ""
		case Token:
			msg = TokenAuth(ctx)
		default:
			msg = CookeAuth(ctx)
		}
		if msg != "" {
			ctx.Abort()
			respx.Exception(ctx, respx.AuthErr, msg)
		} else {
			ctx.Next()
		}
		return
	}
}

func SetAuthType(ctx *gin.Context, authType string) {
	ctx.Request.Header.Set(AuthType, authType)
}

// cookie鉴权
func TokenAuth(ctx *gin.Context) string {
	var token = ctx.Request.Header.Get("Authorization")
	if token == "" {
		return "token is required"
	}
	var user, err = GetUserByToken(token)
	if err != nil || user == nil {
		return "token parse failed"
	}
	if user.GetTokenCache() == "" {
		return "token is expired"
	}
	ctx.Set("user", user)
	return ""
}

// token鉴权
func CookeAuth(ctx *gin.Context) string {
	cookie, err := ctx.Cookie(Cookie)
	if err != nil {
		return "cookie is required"
	}
	var username string
	username, err = encryptx.RSA().Decrypt(cookie)
	if err != nil {
		return "cookie decrypt failed"
	}
	var user = &User{Username: username}
	if token := user.GetTokenCache(); token == "" {
		return "cookie is expired"
	} else {
		user, err = GetUserByToken(token)
		if err != nil {
			return err.Error()
		}
	}
	ctx.Set("user", user)
	return ""
}
