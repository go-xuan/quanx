package authx

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-xuan/quanx/console/redisx"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/go-xuan/quanx/utils/httpx"
)

var SecretKey []byte

// 用户token参数
type User struct {
	Id         string `json:"id"`          // 用户ID
	Account    string `json:"account"`     // 用户名
	Name       string `json:"name"`        // 用户名
	Phone      string `json:"phone"`       // 登录手机
	LoginIp    string `json:"loginIp"`     // 登录IP
	DomainUrl  string `json:"domainUrl"`   // 域名
	ExpireTime int64  `json:"sessionTime"` // 过期时间
}

func (u *User) RedisCacheKey() string {
	return "login@token@" + u.Account
}

// 设置token缓存
func (u *User) SetTokenCache(token string, expiration time.Duration) {
	redisx.GetCmd("user").Set(context.Background(), u.RedisCacheKey(), token, expiration)
}

// 获取token缓存
func (u *User) GetTokenCache() string {
	return redisx.GetCmd("user").Get(context.Background(), u.RedisCacheKey()).Val()
}

func InitSecretKey(key string) {
	SecretKey = []byte(key)
}

// 获取用户ID
func GetUserId(context *gin.Context) (userId int64) {
	tp, err := GetUserByToken(context.Request.Header.Get(httpx.Authorization))
	if err != nil {
		return
	}
	userId, err = strconv.ParseInt(tp.Id, 10, 64)
	if err != nil {
		return
	}
	return
}

// 生成token
func GetTokenByUser(user *User) (token string, err error) {
	var bytes []byte
	bytes, err = json.Marshal(user)
	var mapClaims jwt.MapClaims
	err = json.Unmarshal(bytes, &mapClaims)
	if err != nil {
		return
	}
	token, err = generateToken(mapClaims)
	if err != nil {
		return
	}
	return
}

// 解析token
func GetUserByToken(token string) (user *User, err error) {
	var kvmap map[string]interface{}
	kvmap, err = parseToken(token)
	if err != nil {
		err = errors.New("解析token失败")
		return
	}
	user = &User{
		Id:         kvmap["id"].(string),
		Account:    kvmap["account"].(string),
		Name:       kvmap["name"].(string),
		Phone:      kvmap["phone"].(string),
		LoginIp:    kvmap["loginIp"].(string),
		DomainUrl:  kvmap["domainUrl"].(string),
		ExpireTime: int64(kvmap["sessionTime"].(float64)),
	}
	return
}

// 生成Token值
func generateToken(mapClaims jwt.MapClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims).SignedString(SecretKey)
}

// 解析token
func parseToken(token string) (map[string]interface{}, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	var result = make(map[string]interface{})
	result = claim.Claims.(jwt.MapClaims)
	return result, nil
}
