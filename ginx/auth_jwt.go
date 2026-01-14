package ginx

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-xuan/quanx/cachex"
	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
	"github.com/go-xuan/typex"
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/stringx"
	"github.com/golang-jwt/jwt/v4"
)

// NewJwtValidator 创建JWT验证器
func NewJwtValidator(secret string) *JwtValidator {
	return &JwtValidator{
		Secret: secret,
		White:  make(map[string]string),
	}
}

// JwtValidator JWT验证器
type JwtValidator struct {
	Secret string            `json:"secret" yaml:"secret"` // JWT密钥
	White  map[string]string `json:"white" yaml:"white"`   // 鉴权白名单，map[URL路径]HTTP方法，*表示支持所有方法
	Cache  *cachex.Config    `json:"cache" yaml:"cache"`   // 缓存客户端
}

func (v *JwtValidator) Valid() bool {
	return v.Secret != ""
}

func (v *JwtValidator) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("auth.yaml"),
		configx.NewFileReader("auth.yaml"),
	}
}

func (v *JwtValidator) Execute() error {
	if v.Valid() {
		SetAuthValidator(v)
		// 设置缓存客户端
		if v.Cache != nil && v.Cache.Valid() {
			client, err := cachex.NewClient(v.Cache)
			if err != nil {
				return errorx.Wrap(err, "create cache client failed")
			}
			authCache = client
		}
	}
	return nil
}

// Encrypt 加密用户信息
func (v *JwtValidator) Encrypt(user AuthUser) (string, error) {
	if jwtUser, ok := user.(*JwtUser); ok {
		ciphertext, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtUser).SignedString([]byte(v.Secret))
		if err != nil {
			return "", errorx.Wrap(err, "jwt encrypt token failed")
		}
		return ciphertext, nil
	}
	return "", errorx.New("user is not of type JwtUser")
}

// Decrypt 解密用户信息
func (v *JwtValidator) Decrypt(ciphertext string) (AuthUser, error) {
	var user = &JwtUser{}
	if jt, err := jwt.ParseWithClaims(ciphertext, user, func(*jwt.Token) (interface{}, error) {
		return []byte(v.Secret), nil
	}); err != nil || !jt.Valid {
		return nil, errorx.Wrap(err, "jwt decrypt token failed")
	}
	return user, nil
}

// AddWhite 添加白名单
func (v *JwtValidator) AddWhite(url string, method string) {
	v.White[url] = strings.ToUpper(method)
}

// MatchWhite 匹配白名单
func (v *JwtValidator) MatchWhite(ctx *gin.Context) bool {
	if len(v.White) == 0 {
		return false
	}
	method, url := ctx.Request.Method, ctx.Request.URL.Path
	for u, m := range v.White {
		if stringx.MatchUrl(url, u) && matchMethod(method, strings.ToUpper(m)) {
			return true
		}
	}
	return false
}

func matchMethod(method, rule string) bool {
	return rule == "*" || strings.Contains(rule, method)
}

// Validate 验证用户信息
func (v *JwtValidator) Validate(method AuthMethod) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !v.MatchWhite(ctx) {
			if err := v.validate(ctx, method); err != nil {
				Forbidden(ctx, errorx.Wrap(err, "auth validate failed"))
				ctx.Abort()
				return
			}
		}
	}
}

// validate 验证用户信息
func (v *JwtValidator) validate(ctx *gin.Context, method AuthMethod) error {
	if user := GetSessionUser(ctx); user == nil {
		// 获取鉴权字符串
		if authString, err := GetAuthString(ctx, method); err != nil {
			return errorx.Wrap(err, "get auth string failed")
		} else if user, err = v.Decrypt(authString); err != nil {
			return errorx.Wrap(err, "decrypt token failed")
		} else if !AuthCache().Exist(ctx, user.GetUserId().String()) {
			return errorx.New("auth user is invalid")
		}
		SetSessionUser(ctx, user)
	}
	return nil
}

// Debug 调试模式，模拟当前会话用户，适用于免登录场景
func (v *JwtValidator) Debug(ctx *gin.Context) {
	SetSessionUser(ctx, &JwtUser{Id: 999999999, Name: "DEBUG"})
}

// JwtUser 实现AuthUser
type JwtUser struct {
	Id     int64  `json:"id"`     // 用户id
	Name   string `json:"name"`   // 用户名
	Expire int64  `json:"expire"` // 有效期时间戳
}

// Valid 验证用户信息是否有效，用于实现 jwt.Claims 接口
func (u *JwtUser) Valid() error {
	if u.Id == 0 {
		return errorx.New("user id is empty")
	} else if u.Name == "" {
		return errorx.New("user name is empty")
	} else if u.Expire < time.Now().Unix() {
		return errorx.New("user has expired")
	}
	return nil
}

func (u *JwtUser) GetUserId() typex.Value {
	return typex.NewInt64(u.Id)
}

func (u *JwtUser) GetUsername() string {
	return u.Name
}
