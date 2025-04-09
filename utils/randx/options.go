package randx

import (
	"strconv"
	"strings"

	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/types/timex"
)

const (
	typeInt      = "int"      // 数字
	typeFloat    = "float"    // 浮点数
	typeSequence = "sequence" // 数字编号
	typeTime     = "time"     // 时间
	typeDate     = "date"     // 日期
	typeUUID     = "uuid"     // uuid
	typePhone    = "phone"    // 手机号
	typeName     = "name"     // 姓名
	typeIdCard   = "id_card"  // 身份证
	typePlateNo  = "plate_no" // 车牌号
	typeEmail    = "email"    // 邮箱
	typeIP       = "ip"       // ip地址
	typeProvince = "province" // 省
	typeCity     = "city"     // 市
	typePassword = "password" // 密码
	typeEnum     = "enum"     // 枚举
)

// Options 随机生成
type Options struct {
	Type    string   // 数据类型
	Param   *Param   // 约束条件参数
	Default string   // 默认值
	Offset  int      // 偏移量
	Enums   []string // 枚举
}

// NewRandValue 生成随机数
func (o *Options) NewRandValue() any {
	switch o.Type {
	case typeInt:
		return o.NewInt()
	case typeFloat:
		return o.NewFloat()
	default:
		return o.newString()
	}
}

// NewString 生成随机数
func (o *Options) NewString() string {
	switch o.Type {
	case typeInt:
		return strconv.Itoa(o.NewInt())
	case typeFloat:
		return strconv.FormatFloat(o.NewFloat(), 'f', -1, 64)
	default:
		return o.newString()
	}
}

// NewInt 生成随机数字
func (o *Options) NewInt() int {
	if o.Default == "" {
		return o.Param.Int()
	} else {
		return stringx.ParseInt(o.Default)
	}
}

// NewFloat 生成随机浮点数
func (o *Options) NewFloat() float64 {
	if o.Default == "" {
		return o.Param.Float()
	} else {
		return stringx.ParseFloat(o.Default)
	}
}

// 生成随机字符串
func (o *Options) newString() string {
	var value string
	if param, def := o.Param, o.Default; param != nil && def == "" {
		if o.Type == typeSequence {
			value = strconv.Itoa(stringx.ParseInt(param.Min) + o.Offset)
		} else {
			value = o.randString()
		}
		// 字符替换
		if param.Old != "" && param.New != "" {
			value = strings.ReplaceAll(value, param.Old, param.New)
		}
		// 补充前后缀
		if param.Prefix != "" {
			value = param.Prefix + value
		}
		if param.Suffix != "" {
			value = value + param.Suffix
		}
		// 转换大小写
		if param.Upper {
			value = strings.ToUpper(value)
		}
		if param.Lower {
			value = strings.ToLower(value)
		}
	} else {
		value = def
	}
	return value
}

// 生成随机字符串
func (o *Options) randString() string {
	switch o.Type {
	case typePhone:
		return Phone()
	case typeName:
		return Name()
	case typeIdCard:
		return IdCard()
	case typePlateNo:
		return PlateNo()
	case typeEmail:
		return Email()
	case typeUUID:
		return UUID()
	case typeIP:
		return IP()
	case typeProvince:
		return Province()
	case typeCity:
		return City()
	case typePassword:
		return o.Param.Password()
	case typeDate:
		return o.Param.TimeFmt(timex.DateFmt)
	case typeTime:
		return o.Param.TimeFmt()
	case typeEnum:
		from := append(o.Enums, o.Param.Enums...)
		return StringFrom(from...)
	default:
		return String(o.Param.Length)
	}
}
