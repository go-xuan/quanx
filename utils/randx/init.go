package randx

import (
	"hash/crc32"
	"math/rand"
)

var newRand *rand.Rand

// 初始化
func init() {
	if newRand == nil {
		newRand = rand.New(rand.NewSource(seed()))
	}
}

// 随机种子
func seed() int64 {
	return int64(crc32.ChecksumIEEE([]byte(UUID())))
}

// 随机数据生成类型
type RandType uint

const (
	StringType    RandType = iota // 字符串
	IntType                       // 纯数字
	FloatType                     // 纯浮点
	NoType                        // 编号
	IntStringType                 // 数字字符，可拼接前后缀
	TimeType                      // 时间
	DateType                      // 日期
	UuidType                      // uuid
	PhoneType                     // 手机号
	NameType                      // 姓名
	IdCardType                    // 身份证
	PlateNoType                   // 车牌号
	EmailType                     // 邮箱
	IPType                        // ip地址
	ProvinceType                  // 省
	CityType                      // 市
	PasswordType                  // 密码
	OptionType                    // 枚举
	DatabaseType                  // 数据库取值
)
