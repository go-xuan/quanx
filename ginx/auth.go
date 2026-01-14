package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xuan/typex"
	"github.com/go-xuan/utilx/errorx"

	"github.com/go-xuan/quanx/cachex"
	"github.com/go-xuan/quanx/constx"
)

// AuthMethod 鉴权方式
type AuthMethod string

const (
	CookieAuth    AuthMethod = "cookie"        // cookie方式鉴权
	TokenAuth     AuthMethod = "token"         // token方式鉴权
	cookieAuthKey            = "COOKIE_USER"   // cookie获取键
	tokenAuthKey             = "Authorization" // token获取键
)

var (
	authValidator AuthValidator // 鉴权验证器
	authCache     cachex.Client // auth缓存客户端
)

// AuthValidate 获取鉴权验证器
func AuthValidate() AuthValidator {
	if authValidator == nil {
		authValidator = NewJwtValidator(constx.DefaultSource)
	}
	return authValidator
}

// AuthCache 获取auth缓存
func AuthCache() cachex.Client {
	if authCache == nil {
		authCache = cachex.GetClient("auth")
	}
	return authCache
}

// SetAuthValidator 自定义鉴权验证器
func SetAuthValidator(validator AuthValidator) {
	if validator != nil {
		authValidator = validator
	}
}

// AuthValidator 鉴权验证器
type AuthValidator interface {
	AddWhite(url, method string)                // 添加白名单
	Encrypt(user AuthUser) (string, error)      // 加密用户信息
	Decrypt(auth string) (AuthUser, error)      // 解密用户信息
	Validate(method AuthMethod) gin.HandlerFunc // 鉴权中间件
}

// AuthUser 鉴权用户
type AuthUser interface {
	GetUserId() typex.Value // 用户id
	GetUsername() string    // 用户名
}

// GetAuthString 获取鉴权字符串
func GetAuthString(ctx *gin.Context, method AuthMethod) (string, error) {
	switch method {
	case CookieAuth:
		cookie, err := ctx.Cookie(cookieAuthKey)
		if err != nil {
			return "", errorx.Wrap(err, "get cookie failed")
		} else if cookie == "" {
			return "", errorx.New("auth cookie is empty")
		}
		return cookie, nil
	case TokenAuth:
		token := ctx.Request.Header.Get(tokenAuthKey)
		if token == "" {
			return "", errorx.New("auth token is empty")
		}
		return token, nil
	default:
		return "", errorx.New("auth method not support")
	}
}

// SetSessionUser 设置会话用户
func SetSessionUser(ctx *gin.Context, user AuthUser) {
	ctx.Set(sessionUserKey, user)
}

// GetSessionUser 获取会话用户
func GetSessionUser(ctx *gin.Context) AuthUser {
	if value, exist := ctx.Get(sessionUserKey); exist {
		if user, ok := value.(AuthUser); ok {
			return user
		}
	}
	return nil
}

// SetAuthCookie 设置身份验证cookie
func SetAuthCookie(ctx *gin.Context, cookie string, maxAge int) {
	ctx.SetCookie(cookieAuthKey, cookie, maxAge, "", "", false, true)
}

// RemoveAuthCookie 移除身份验证cookie maxAge=-1即可移除cookie
func RemoveAuthCookie(ctx *gin.Context) {
	ctx.SetCookie(cookieAuthKey, "", -1, "", "", false, true)
}
