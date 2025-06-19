package ginx

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/extra/cachex"
	"github.com/go-xuan/quanx/types/anyx"
)

var (
	authValidator   AuthValidator // 鉴权验证器
	authCacheClient cachex.Client // token缓存客户端
)

// SetAuthValidator 设置鉴权验证器
func SetAuthValidator(validator AuthValidator) {
	if validator != nil {
		authValidator = validator
	}
}

// AuthValidate 身份验证器
func AuthValidate() AuthValidator {
	if authValidator == nil {
		authValidator = NewJwtValidator("123456")
	}
	return authValidator
}

func AuthCache() cachex.Client {
	if authCacheClient == nil {
		authCacheClient = cachex.GetClient("auth")
	}
	return authCacheClient
}

// AuthValidator 验证器
type AuthValidator interface {
	Encrypt(user AuthUser) (string, error) // 用户信息加密
	TokenValidate(ctx *gin.Context)        // token方式鉴权
	CookieValidate(ctx *gin.Context)       // cookie方式鉴权
	Ignore(ignores ...string)              // 添加白名单
	IsIgnore(ctx *gin.Context) bool        // 白名单免鉴权
}

// AuthUser 鉴权用户
type AuthUser interface {
	UserId() anyx.Value // 用户ID
	Username() string   // 用户名
	TTL() int           // 存活时长，秒
}

// SetSessionUser 设置会话用户
func SetSessionUser(ctx *gin.Context, user AuthUser) {
	ctx.Set(sessionUserKey, user)
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
func SetAuthCookie(ctx *gin.Context, user AuthUser) (string, error) {
	if cookie, err := AuthValidate().Encrypt(user); err != nil {
		return "", errorx.Wrap(err, "new cookie error")
	} else if err = AuthCache().Set(ctx, user.UserId().String(), cookie, time.Duration(user.TTL())*time.Second); err != nil {
		return "", errorx.Wrap(err, "save cookie to cache error")
	} else {
		ctx.SetCookie(userCookieKey, cookie, user.TTL(), "", "", false, true)
		SetSessionUser(ctx, user)
		return cookie, nil
	}
}

// RemoveAuthCookie 移除身份验证cookie maxAge=-1即可移除cookie
func RemoveAuthCookie(ctx *gin.Context) {
	ctx.SetCookie(userCookieKey, "", -1, "", "", false, true)
}

// SetAuthToken 生成并缓存token
func SetAuthToken(ctx *gin.Context, user AuthUser) (string, error) {
	if token, err := AuthValidate().Encrypt(user); err != nil {
		return "", errorx.Wrap(err, "new token error")
	} else if err = AuthCache().Set(ctx, user.UserId().String(), token, time.Duration(user.TTL())*time.Second); err != nil {
		return "", errorx.Wrap(err, "save token to cache error")
	} else {
		SetSessionUser(ctx, user)
		return token, nil
	}
}

// RemoveAuthToken 移除token
func RemoveAuthToken(ctx *gin.Context, userId string) {
	AuthCache().Delete(ctx, userId)
}
