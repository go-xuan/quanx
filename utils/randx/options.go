package randx

import (
	"strings"
	"time"

	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/types/timex"
)

const (
	IntType       = "int"        // 数字
	FloatType     = "float"      // 浮点数
	SequenceType  = "sequence"   // 数字编号
	TimeType      = "time"       // 时间
	DateType      = "date"       // 日期
	UuidType      = "uuid"       // uuid
	PhoneType     = "phone"      // 手机号
	NameType      = "name"       // 姓名
	IdCardType    = "id_card"    // 身份证
	PlateNoType   = "plate_no"   // 车牌号
	EmailType     = "email"      // 邮箱
	IPType        = "ip"         // ip地址
	ProvinceType  = "province"   // 省
	CityType      = "city"       // 市
	PasswordType  = "password"   // 密码
	EnumType      = "enum"       // 枚举
	OtherEnumType = "other_enum" // 其他枚举
)

// 随机生成
type Options struct {
	Type    string   // 数据类型
	Param   *Param   // 约束条件参数
	Default string   // 默认值
	Offset  int      // 偏移量
	Enums   []string // 枚举
}

// 生成随机数
func (o *Options) RandData() any {
	switch o.Type {
	case IntType:
		return o.RandInt()
	case FloatType:
		return o.RandFloat()
	default:
		return o.RandString()
	}
}

// 生成随机数
func (o *Options) RandDataString() string {
	switch o.Type {
	case IntType:
		return stringx.ParseInt(o.RandInt())
	case FloatType:
		return stringx.ParseFloat64(o.RandFloat())
	default:
		return o.RandString()
	}
}

// 生成随机数字
func (o *Options) RandInt() int {
	if o.Default == "" {
		return o.Param.Int()
	} else {
		return stringx.ToInt(o.Default)
	}
}

// 生成随机浮点数
func (o *Options) RandFloat() float64 {
	if o.Default == "" {
		return o.Param.Float()
	} else {
		return stringx.ToFloat64(o.Default)
	}
}

// 生成随机字符串
func (o *Options) RandString() (result string) {
	if param, def := o.Param, o.Default; param != nil && def == "" {
		if o.Type == SequenceType {
			result = stringx.ParseInt(stringx.ToInt(param.Min) + o.Offset)
		} else {
			result = o.randString()
		}
		// 字符替换
		if param.Old != "" && param.New != "" {
			result = strings.ReplaceAll(result, param.Old, param.New)
		}
		// 补充前后缀
		if param.Prefix != "" {
			result = param.Prefix + result
		}
		if param.Suffix != "" {
			result = result + param.Suffix
		}
		// 转换大小写
		if param.Upper {
			result = strings.ToUpper(result)
		}
		if param.Lower {
			result = strings.ToLower(result)
		}
	} else {
		result = def
	}
	return
}

// 生成随机字符串
func (o *Options) randString() (result string) {
	switch o.Type {
	case PhoneType:
		return Phone()
	case NameType:
		return Name()
	case IdCardType:
		return IdCard()
	case PlateNoType:
		return PlateNo()
	case EmailType:
		return Email()
	case UuidType:
		return UUID()
	case IPType:
		return IP()
	case ProvinceType:
		return Province()
	case CityType:
		return City()
	case PasswordType:
		return o.Param.Password()
	case DateType:
		return o.Param.TimeFmt(timex.DateFmt)
	case TimeType:
		return o.Param.TimeFmt()
	case EnumType:
		return Enum(o.Param.Enums)
	case OtherEnumType:
		return Enum(append(o.Enums, o.Param.Enums...))
	default:
		return String(o.Param.Length)
	}
}

// 随机数生成参数
type Param struct {
	Min    string   // 最小值
	Max    string   // 最大值
	Prefix string   // 前缀
	Suffix string   // 后缀
	Upper  bool     // 转大写
	Lower  bool     // 转小写
	Old    string   // 替换旧字符
	New    string   // 替换新字符
	Format string   // 时间格式
	Length int      // 长度
	Prec   int      // 小数位精度
	Level  int      // 级别
	Enums  []string // 枚举选项，多个以逗号分割
}

func NewParam(constraint string) *Param {
	params := stringx.ParseUrlParams(constraint)
	return &Param{
		Min:    params["min"],
		Max:    params["max"],
		Prefix: params["prefix"],
		Suffix: params["suffix"],
		Upper:  params["upper"] == "true",
		Lower:  params["suffix"] == "true",
		Old:    params["old"],
		New:    params["new"],
		Format: params["format"],
		Length: stringx.ToInt(params["length"]),
		Prec:   stringx.ToInt(params["prec"]),
		Level:  stringx.ToInt(params["level"]),
		Enums:  strings.Split(params["enums"], ","),
	}
}

func (c *Param) Password() string {
	switch c.Level {
	case 2:
		return Password(WithNumber|WithLowerLetter|WithUpperLetter, c.Length)
	case 3:
		return Password(WithNumber|WithLowerLetter|WithUpperLetter|WithSpecial, c.Length)
	default:
		return Password(WithNumber, c.Length)
	}
}

func (c *Param) Int() int {
	var minv, maxv int
	if minv = stringx.ToInt(c.Min); minv == 0 {
		minv = 1
	}
	if maxv = stringx.ToInt(c.Max); maxv == 0 {
		maxv = 999
	}
	return IntRange(minv, maxv)
}

func (c *Param) Float() float64 {
	var minv, maxv float64
	var prec int
	if minv = stringx.ToFloat64(c.Min); minv == 0 {
		minv = 1
	}
	if maxv = stringx.ToFloat64(c.Max); maxv == 0 {
		maxv = 999
	}
	if prec = c.Prec; prec == 0 {
		prec = 6
	}
	return Float64Range(minv, maxv, prec)
}

func (c *Param) Time() time.Time {
	end := time.Now()
	if c.Max != "" {
		end = timex.ToTime(c.Max)
	}
	start := end.Add(time.Hour * -24 * 30)
	if c.Min != "" {
		start = timex.ToTime(c.Min)
	}
	return TimeRange(start, start)
}

func (c *Param) TimeFmt(layouts ...string) string {
	end := time.Now()
	if c.Max != "" {
		end = timex.ToTime(c.Max)
	}
	start := end.Add(time.Hour * -24 * 30)
	if c.Min != "" {
		start = timex.ToTime(c.Min)
	}
	var layout = timex.TimeFmt
	if len(layouts) > 0 {
		layout = layouts[0]
	} else if c.Format != "" {
		layout = c.Format
	}
	return TimeRange(start, start).Format(layout)
}

func (c *Param) Enum() string {
	return Enum(c.Enums)
}

func (c *Param) String() string {
	return String(c.Length)
}
