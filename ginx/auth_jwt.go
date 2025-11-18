package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xuan/typex"
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/stringx"
	"github.com/golang-jwt/jwt/v4"
)

// NewJwtValidator 创建JWT验证器
func NewJwtValidator(secret string, ignore ...string) *JwtValidator {
	return &JwtValidator{
		secret: []byte(secret),
		ignore: ignore,
	}
}

// JwtValidator JWT验证器
type JwtValidator struct {
	secret []byte   // jwt密钥
	ignore []string // 鉴权白名单
}

// Encrypt 加密用户信息
func (v *JwtValidator) Encrypt(user AuthUser) (string, error) {
	if jwtUser, ok := user.(*JwtUser); ok {
		ciphertext, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtUser).SignedString(v.secret)
		if err != nil {
			return "", errorx.Wrap(err, "jwt encrypt token error")
		}
		return ciphertext, nil
	}
	return "", errorx.New("user is not of type JwtUser")
}

// Decrypt 解密用户信息
func (v *JwtValidator) Decrypt(ciphertext string) (AuthUser, error) {
	var user = &JwtUser{}
	if jt, err := jwt.ParseWithClaims(ciphertext, user, func(*jwt.Token) (interface{}, error) {
		return v.secret, nil
	}); err != nil || !jt.Valid {
		return nil, errorx.Wrap(err, "jwt decrypt error")
	}
	return user, nil
}

// AddWhite 添加白名单
func (v *JwtValidator) AddWhite(url ...string) {
	v.ignore = append(v.ignore, url...)
}

// MatchWhite 匹配白名单
func (v *JwtValidator) MatchWhite(ctx *gin.Context) bool {
	url := ctx.Request.URL.Path
	for _, ignore := range v.ignore {
		if stringx.MatchUrl(url, ignore) {
			return true
		}
	}
	return false
}

// Validate 验证用户信息
func (v *JwtValidator) Validate(method AuthMethod) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !v.MatchWhite(ctx) {
			user, err := v.validate(ctx, method)
			if err != nil {
				Forbidden(ctx, errorx.Wrap(err, "auth validate failed"))
				ctx.Abort()
				return
			}
			SetSessionUser(ctx, user)
		}
	}
}

func (v *JwtValidator) validate(ctx *gin.Context, method AuthMethod) (AuthUser, error) {
	// 如果session中存在用户，直接返回，避免重复解密
	var user AuthUser
	if user = GetSessionUser(ctx); user != nil {
		return user, nil
	}
	// 获取鉴权字符串
	authString, err := GetAuthString(ctx, method)
	if err != nil {
		return nil, errorx.Wrap(err, "get auth string error")
	}
	// 从token中解密用户信息
	if user, err = v.Decrypt(authString); err != nil {
		return nil, errorx.Wrap(err, "decrypt error")
	}
	// 验证用户信息是否有效
	if !AuthCache().Exist(ctx, user.GetUserId().String()) {
		return nil, errorx.Wrap(err, "auth user is invalid")
	}
	return user, nil
}

// Debug 调试模式，模拟当前会话用户，适用于免登录场景
func (v *JwtValidator) Debug(ctx *gin.Context) {
	SetSessionUser(ctx, &JwtUser{Id: 999999999, Name: "DEBUG"})
}

// JwtUser 实现AuthUser
type JwtUser struct {
	Id   int64  `json:"id"`   // 用户ID
	Name string `json:"name"` // 用户名
}

// Valid 验证用户信息是否有效，用于实现 jwt.Claims 接口
func (u *JwtUser) Valid() error {
	if u.Id == 0 {
		return errorx.New("user id is empty")
	} else if u.Name == "" {
		return errorx.New("user name is empty")
	}
	return nil
}

func (u *JwtUser) GetUserId() typex.Value {
	return typex.NewInt64(u.Id)
}

func (u *JwtUser) GetUsername() string {
	return u.Name
}
