package ginx

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/extra/cachex"
	"github.com/go-xuan/quanx/types/anyx"
)

// AuthValidate 默认使用jwt验证
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
	Encrypt() (string, error)   // 用户信息加密
	Decrypt(token string) error // 用户信息解密
	UserId() anyx.Value         // 用户ID
	TTL() time.Duration         // 存活时间
}

var authCacheClient cachex.Client // token缓存客户端

func AuthCache() cachex.Client {
	if authCacheClient == nil {
		authCacheClient = cachex.GetClient("auth")
	}
	return authCacheClient
}

// SetSessionUser 设置会话用户
func SetSessionUser(ctx *gin.Context, user AuthUser) {
	ctx.Set(sessionUserKey, user)
	// token续命
	_ = AuthCache().Expire(ctx, user.UserId().String(), user.TTL())
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
func SetAuthCookie(ctx *gin.Context, user AuthUser) error {
	if cookie, err := user.Encrypt(); err != nil {
		return errorx.Wrap(err, "new cookie error")
	} else {
		ctx.SetCookie(userCookieKey, cookie, int(user.TTL().Seconds()), "", ctx.Request.Host, false, true)
	}
	return nil
}

// RemoveAuthCookie 移除身份验证cookie maxAge=-1即可移除cookie
func RemoveAuthCookie(ctx *gin.Context) {
	ctx.SetCookie(userCookieKey, "", -1, "", "", false, true)
}

// ValidateCookie cookie鉴权
func ValidateCookie(ctx *gin.Context, user AuthUser) error {
	var cookie, err = ctx.Cookie(userCookieKey)
	if err != nil || cookie == "" {
		return errorx.Wrap(err, "get request cookie failed")
	}
	if err = user.Decrypt(cookie); err != nil {
		return errorx.Wrap(err, "cookie parse failed")
	}
	if !AuthCache().Exist(ctx, user.UserId().String()) {
		return errorx.New("cookie has expired")
	}
	SetSessionUser(ctx, user)
	return nil
}

// SetToken 生成并缓存token
func SetToken(ctx *gin.Context, user AuthUser) (string, error) {
	if token, err := user.Encrypt(); err != nil {
		return "", errorx.Wrap(err, "new token error")
	} else if err = AuthCache().Set(ctx, user.UserId().String(), token, user.TTL()); err != nil {
		return "", errorx.Wrap(err, "save token to cache error")
	} else {
		return token, nil
	}
}

// RemoveToken 移除token
func RemoveToken(ctx *gin.Context, userId string) {
	AuthCache().Delete(ctx, userId)
}

// ValidateToken token鉴权
func ValidateToken(ctx *gin.Context, user AuthUser) error {
	var token = ctx.Request.Header.Get("Authorization")
	if token == "" {
		return errorx.New("get request token failed")
	}
	if err := user.Decrypt(token); err != nil {
		return errorx.Wrap(err, "token parse failed")
	}
	if !AuthCache().Exist(ctx, user.UserId().String()) {
		return errorx.New("token has expired")
	}
	SetSessionUser(ctx, user)
	return nil
}
