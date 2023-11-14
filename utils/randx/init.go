package randx

import (
	"hash/crc32"
	"math/rand"
)

// 随机数据生成约束条件关键字
const (
	Min       = "min"        // 最小值
	Max       = "max"        // 最大值
	Prefix    = "prefix"     // 前缀
	Suffix    = "suffix"     // 后缀
	Prec      = "prec"       // 小数位
	Format    = "format"     // 时间格式
	Length    = "length"     // 长度
	Lower     = "lower"      // 小写
	Upper     = "upper"      // 大写
	HasNumber = "has_number" // 是否含有数字
	HasSymbol = "has_symbol" // 是否含有符号
	Old       = "old"        // 替换旧字符
	New       = "new"        // 替换新字符
	Options   = "options"    // 枚举选项，多个以逗号分割
	Table     = "table"      // 表名，表字段所属表名
	Field     = "field"      // 字段名
)

// 字典配置
const (
	LetterLower    = "abcdefghijklmnopqrstuvwxyz"
	LetterUpper    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharLower      = "abcdefghijklmnopqrstuvwxyz1234567890"
	CharAll        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	NUMBER         = "1234567890"
	SPECIAL        = "!@#$%&*+-=?"
	ProvinceSimple = "京,津,冀,晋,蒙,辽,吉,黑,沪,苏,浙,皖,闽,赣,鲁,豫,鄂,湘,粤,桂,琼,渝,川,贵,云,藏,陕,甘,青,宁,新,台,港,澳"
	ProvinceName   = "北京市,天津市,河北省,山西省,内蒙古自治区,辽宁省,吉林省,黑龙江省,上海市,江苏省,浙江省,安徽省,福建省,江西省,山东省,河南省,湖北省,湖南省,广东省,广西壮族自治区,海南省,重庆市,四川省,贵州省,云南省,西藏自治区,陕西省,甘肃省,青海省,宁夏回族自治区,新疆维吾尔自治区,台湾省,香港特别行政区,澳门特别行政区"
	HubeiCityName  = "武汉,黄石,十堰,宜昌,襄阳,鄂州,荆门,孝感,荆州,黄冈,咸宁,随州,恩施,仙桃,潜江,天门,神农架"
	HubeiPostcode  = "420100,420200,420300,420500,420600,420700,420800,420900,421000,421100,421200,421300,422800,429004,429005,429006,429021"
	FamilyNameCn   = "赵,钱,孙,李,周,吴,黄,高,郑,王,冯,陈,蒋,沈,韩,杨,朱,秦,许,何,吕,张,孔,曹"
	NumberCn       = "一,二,三,四,五,六,七,八,九,十,百,千,万,亿"
	ShengXiao      = "鼠,牛,虎,兔,龙,蛇,马,羊,猴,鸡,狗,猪"
	PhonePrefix    = "358"
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
