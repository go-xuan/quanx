package ginx

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// jwt秘钥
var SecretKey []byte

// 初始化jwt秘钥
func SetSecretKey(key []byte) {
	SecretKey = key
}

// 用户token参数
type User struct {
	Id         int64  `json:"id"`         // 用户ID
	Account    string `json:"account"`    // 用户账号
	Name       string `json:"name"`       // 用户姓名
	Phone      string `json:"phone"`      // 登录手机
	Ip         string `json:"ip"`         // 登录IP
	Domain     string `json:"domain"`     // 域名
	ExpireTime int64  `json:"expireTime"` // 过期时间
}

// 获取用户ID
func GetUserId(context *gin.Context) (userId int64) {
	var err error
	var userData = make(map[string]interface{})
	if userData, err = parseToken(context.Request.Header.Get(Authorization)); err != nil {
		return
	}
	userId = userData["id"].(int64)
	return
}

// 生成token
func NewToken(user *User) (token string, err error) {
	var bytes []byte
	if bytes, err = json.Marshal(user); err != nil {
		return
	}
	var mapClaims jwt.MapClaims
	if err = json.Unmarshal(bytes, &mapClaims); err != nil {
		return
	}
	if token, err = generateToken(mapClaims); err != nil {
		return
	}
	return
}

// 解析token
func ParseUserFromToken(token string) (user *User, err error) {
	var userData = make(map[string]interface{})
	if userData, err = parseToken(token); err != nil {
		return
	}
	user = &User{
		Id:         userData["id"].(int64),
		Account:    userData["account"].(string),
		Name:       userData["name"].(string),
		Phone:      userData["phone"].(string),
		Ip:         userData["ip"].(string),
		Domain:     userData["domain"].(string),
		ExpireTime: int64(userData["expireTime"].(float64)),
	}
	return
}

// 生成Token值
func generateToken(mapClaims jwt.MapClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims).SignedString(SecretKey)
}

// 解析token
func parseToken(token string) (map[string]interface{}, error) {
	if claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	}); err != nil {
		return nil, err
	} else {
		var result = make(map[string]interface{})
		result = claim.Claims.(jwt.MapClaims)
		return result, nil
	}
}
