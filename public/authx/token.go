package authx

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/go-xuan/quanx/common/constx"
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
	return RedisKeyPrefix + u.Account
}

func InitSecretKey(key string) {
	SecretKey = []byte(key)
}

// 生成token
func BuildAuthToken(user *User) (string, error) {
	bytes, err := json.Marshal(user)
	var mapClaims jwt.MapClaims
	err = json.Unmarshal(bytes, &mapClaims)
	if err != nil {
		return "", err
	}
	return GenerateToken(mapClaims)
}

// 解析token
func ParseAuthToken(token string) (user *User, err error) {
	var resultMap map[string]interface{}
	resultMap, err = ParseTokenToJson(token)
	if err != nil {
		err = errors.New("解析token失败")
		return
	}
	user = &User{
		Id:         resultMap["id"].(string),
		Account:    resultMap["account"].(string),
		Name:       resultMap["name"].(string),
		Phone:      resultMap["phone"].(string),
		LoginIp:    resultMap["loginIp"].(string),
		DomainUrl:  resultMap["domainUrl"].(string),
		ExpireTime: int64(resultMap["sessionTime"].(float64)),
	}
	return
}

// 获取用户ID
func GetUserId(context *gin.Context) (userId int64) {
	tp, err := ParseAuthToken(context.Request.Header.Get(constx.Authorization))
	if err != nil {
		return
	}
	userId, err = strconv.ParseInt(tp.Id, 10, 64)
	if err != nil {
		return
	}
	return
}

// 生成Token值
func GenerateToken(mapClaims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString(SecretKey)
}

// 解析token
func ParseTokenToJson(token string) (map[string]interface{}, error) {
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
