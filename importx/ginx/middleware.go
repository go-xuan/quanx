package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xuan/quanx/importx/encryptx"
	"github.com/go-xuan/quanx/utilx/anyx"

	"github.com/go-xuan/quanx/commonx/respx"
	"github.com/go-xuan/quanx/utilx/ipx"
)

const (
	AuthType      = "auth_type"     // 鉴权方式
	NoAuth        = "no"            // 免鉴权标识
	Token         = "token"         // token标识
	Cookie        = "cookie"        // cookie鉴权标识
	Authorization = "Authorization" // cookie鉴权标识
)

func SetAuthType(ctx *gin.Context, authType string) {
	ctx.Request.Header.Set(AuthType, authType)
}

func RemoveCookie(ctx *gin.Context) {
	ctx.SetCookie(Cookie, "", -1, "", "", false, true)
}

func SetCookie(ctx *gin.Context, username string, age ...int) {
	if cookie, err := encryptx.RSA().Encrypt(username); err != nil {
		ctx.Abort()
		respx.BuildError(ctx, err)
		return
	} else {
		var maxAge = anyx.IfElseValue(len(age) > 0, age[0], 3600)
		ctx.SetCookie(Cookie, cookie, maxAge, "", "", false, true)
	}
}

func CorrectIP(ctx *gin.Context) {
	var ip string
	if ip = ctx.ClientIP(); ip == "::1" {
		ip = ipx.GetWLANIP()
	}
	ctx.Set("ip", ip)
	ctx.Next()
	return
}

// 白名单
func NotAuth(ctx *gin.Context) {
	SetAuthType(ctx, NoAuth)
	ctx.Next()
	return
}

// cookie鉴权
func Auth(ctx *gin.Context) {
	if err := Authenticate(ctx); err != "" {
		ctx.Abort()
		respx.Exception(ctx, respx.AuthErr, err)
	} else {
		ctx.Next()
	}
	return
}

func Authenticate(ctx *gin.Context) string {
	switch ctx.Request.Header.Get(AuthType) {
	case NoAuth:
		return ""
	case Token:
		return TokenAuthenticate(ctx)
	default:
		return CookeAuthenticate(ctx)
	}
}

// cookie鉴权
func TokenAuthenticate(ctx *gin.Context) string {
	if token := ctx.Request.Header.Get(Authorization); token == "" {
		return respx.AuthRequiredErr.Msg
	} else {
		if user, err := GetUserByToken(token); err != nil || user == nil {
			return respx.AuthInvalidErr.Msg
		} else if user.GetTokenCache() == "" {
			return respx.AuthExpiredErr.Msg
		} else {
			ctx.Set("user", user)
		}
		return ""
	}
}

// token鉴权
func CookeAuthenticate(ctx *gin.Context) string {
	if cookie, err := ctx.Cookie(Cookie); err != nil {
		return respx.AuthRequiredErr.Msg
	} else {
		var account string
		if account, err = encryptx.RSA().Decrypt(cookie); err != nil {
			return respx.AuthInvalidErr.Msg
		}
		var user = &User{Account: account}
		if token := user.GetTokenCache(); token == "" {
			return respx.AuthExpiredErr.Msg
		} else if user, err = GetUserByToken(token); err != nil {
			return respx.AuthInvalidErr.Msg
		} else {
			ctx.Set("user", user)
		}
		return ""
	}
}
