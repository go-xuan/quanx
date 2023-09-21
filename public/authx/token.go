package authx

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

// 用户token参数
type User struct {
	UserId     string `json:"userId"`      // 用户ID
	UserName   string `json:"userName"`    // 用户名
	Phone      string `json:"phone"`       // 登录手机
	LoginIp    string `json:"loginIp"`     // 登录IP
	DomainUrl  string `json:"domainUrl"`   // 域名
	ExpireTime int64  `json:"sessionTime"` // 过期时间
}

// 生成token
func BuildAuthToken(user *User, Key string) (string, error) {
	bytes, err := json.Marshal(user)
	var mapClaims jwt.MapClaims
	err = json.Unmarshal(bytes, &mapClaims)
	if err != nil {
		return "", err
	}
	return GenerateToken(mapClaims, Key)
}

// 解析token
func ParseAuthToken(token, key string) (user *User, err error) {
	var resultMap map[string]interface{}
	resultMap, err = ParseTokenToJson(token, key)
	if err != nil {
		log.Error("解析token失败", err)
		return nil, err
	}
	user = &User{
		UserId:     resultMap["userId"].(string),
		UserName:   resultMap["userName"].(string),
		Phone:      resultMap["phone"].(string),
		LoginIp:    resultMap["loginIp"].(string),
		DomainUrl:  resultMap["domainUrl"].(string),
		ExpireTime: int64(resultMap["sessionTime"].(float64)),
	}
	return
}

// 获取用户ID
func GetUserId(context *gin.Context, key string) (userId int64) {
	auth := context.Request.Header.Get("Authorization")
	tp, err := ParseAuthToken(auth, key)
	if err != nil {
		return
	}
	userId, err = strconv.ParseInt(tp.UserId, 10, 64)
	if err != nil {
		return
	}
	return
}

// 生成Token值
func GenerateToken(mapClaims jwt.MapClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString([]byte(secret))
}

// 解析token
func ParseTokenToJson(token string, secret string) (map[string]interface{}, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	var result = make(map[string]interface{})
	result = claim.Claims.(jwt.MapClaims)
	return result, nil
}
