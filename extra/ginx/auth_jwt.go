package ginx

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/respx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
)

func NewJwtValidator(secret string, ignore ...string) *JwtValidator {
	return &JwtValidator{
		secret: []byte(secret),
		ignore: ignore,
	}
}

type JwtValidator struct {
	secret []byte   // 密钥
	ignore []string // 白名单
}

func (v *JwtValidator) Encrypt(user AuthUser) (string, error) {
	if jwtUser, ok := user.(*JwtUser); ok {
		if ciphertext, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtUser).SignedString(v.secret); err != nil {
			return "", errorx.Wrap(err, "jwt encrypt error")
		} else {
			return ciphertext, nil
		}
	} else {
		return "", errorx.New("user is not of type JwtUser")
	}
}

func (v *JwtValidator) TokenValidate(ctx *gin.Context) {
	if !v.Ignore(ctx) {
		if err := v.tokenValidate(ctx); err != nil {
			respx.Forbidden(ctx, errorx.Wrap(err, "auth validate failed"))
			ctx.Abort()
		}
	}
}

func (v *JwtValidator) CookieValidate(ctx *gin.Context) {
	if !v.Ignore(ctx) {
		if err := v.cookieValidate(ctx); err != nil {
			respx.Forbidden(ctx, errorx.Wrap(err, "auth validate failed"))
			ctx.Abort()
		}
	}
}

func (v *JwtValidator) AddWhite(ignore ...string) {
	v.ignore = append(v.ignore, ignore...)
}

func (v *JwtValidator) Ignore(ctx *gin.Context) bool {
	url := ctx.Request.URL.Path
	for _, ignore := range v.ignore {
		if stringx.MatchUrl(url, ignore) {
			return true
		}
	}
	return false
}

func (v *JwtValidator) tokenValidate(ctx *gin.Context) error {
	if token := ctx.Request.Header.Get("Authorization"); token != "" {
		return errorx.New("request token required")
	} else if user, err := v.validate(ctx, token); err != nil {
		return errorx.Wrap(err, "validate cookie error")
	} else if _, err = SetAuthToken(ctx, user); err != nil {
		return errorx.Wrap(err, "set cookie error")
	}
	return nil
}

func (v *JwtValidator) cookieValidate(ctx *gin.Context) error {
	if cookie, _ := ctx.Cookie(userCookieKey); cookie == "" {
		return errorx.New("request cookie required")
	} else if user, err := v.validate(ctx, cookie); err != nil {
		return errorx.Wrap(err, "validate cookie error")
	} else if _, err = SetAuthCookie(ctx, user); err != nil {
		return errorx.Wrap(err, "set cookie error")
	}
	return nil
}

func (v *JwtValidator) validate(ctx *gin.Context, ciphertext string) (*JwtUser, error) {
	var user = &JwtUser{}
	if jt, err := jwt.ParseWithClaims(ciphertext, user, func(*jwt.Token) (interface{}, error) {
		return v.secret, nil
	}); err != nil || !jt.Valid {
		return nil, errorx.Wrap(err, "jwt decrypt error")
	}
	if !AuthCache().Exist(ctx, user.UserId().String()) {
		return nil, errorx.New("auth cache has expired:" + user.UserId().String())
	}
	user.Update = time.Now().Unix()
	return user, nil
}

func (v *JwtValidator) Debug(ctx *gin.Context) {
	SetSessionUser(ctx, &JwtUser{
		Id:     999999999,
		Name:   "debug",
		Update: time.Now().Unix(),
	})
}

// JwtUser 实现AuthUser
type JwtUser struct {
	Id     int64  `json:"id"`     // 用户ID
	Name   string `json:"name"`   // 用户名
	Age    int    `json:"age"`    // 存活时长
	Update int64  `json:"update"` // 更新时间
}

func (u *JwtUser) Valid() error {
	if u.Update+int64(u.Age) < time.Now().Unix() {
		return errorx.New("claims valid failed: current user expired")
	}
	return nil
}

func (u *JwtUser) UserId() anyx.Value {
	return anyx.Int64Value(u.Id)
}

func (u *JwtUser) Username() string {
	return u.Name
}

func (u *JwtUser) TTL() int {
	return anyx.IfZero(u.Age, 3600)
}
