package randx

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-xuan/quanx/types/intx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/types/timex"
)

const (
	typeInt       = "int"        // 数字
	typeFloat     = "float"      // 浮点数
	typeSequence  = "sequence"   // 数字编号
	typeTime      = "time"       // 时间
	typeDate      = "date"       // 日期
	typeUUID      = "uuid"       // uuid
	typePhone     = "phone"      // 手机号
	typeName      = "name"       // 姓名
	typeIdCard    = "id_card"    // 身份证
	typePlateNo   = "plate_no"   // 车牌号
	typeEmail     = "email"      // 邮箱
	typeIP        = "ip"         // ip地址
	typeProvince  = "province"   // 省
	typeCity      = "city"       // 市
	typePassword  = "password"   // 密码
	typeEnum      = "enum"       // 枚举
	typeOtherEnum = "other_enum" // 其他枚举
)

// Options 随机生成
type Options struct {
	Type    string   // 数据类型
	Args    *Args    // 约束条件参数
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
		return o.Args.Int()
	} else {
		return stringx.ParseInt(o.Default)
	}
}

// NewFloat 生成随机浮点数
func (o *Options) NewFloat() float64 {
	if o.Default == "" {
		return o.Args.Float()
	} else {
		return stringx.ParseFloat(o.Default)
	}
}

// 生成随机字符串
func (o *Options) newString() string {
	var value string
	if param, def := o.Args, o.Default; param != nil && def == "" {
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
		return o.Args.Password()
	case typeDate:
		return o.Args.TimeFmt(timex.DateFmt)
	case typeTime:
		return o.Args.TimeFmt()
	case typeEnum:
		return Enum(o.Args.Enums)
	case typeOtherEnum:
		return Enum(append(o.Enums, o.Args.Enums...))
	default:
		return String(o.Args.Length)
	}
}

// Args 随机数生成参数
type Args struct {
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

func NewArgs(args string) *Args {
	params := stringx.ParseUrlParams(args)
	return &Args{
		Min:    params["min"],
		Max:    params["max"],
		Prefix: params["prefix"],
		Suffix: params["suffix"],
		Upper:  params["upper"] == "true",
		Lower:  params["lower"] == "true",
		Old:    params["old"],
		New:    params["new"],
		Format: params["format"],
		Length: stringx.ParseInt(params["length"]),
		Prec:   stringx.ParseInt(params["prec"]),
		Level:  stringx.ParseInt(params["level"]),
		Enums:  strings.Split(params["enums"], ","),
	}
}

func (c *Args) Password() string {
	switch c.Level {
	case 2:
		return StringWith(WithNumber|WithLowerLetter|WithUpperLetter, c.Length)
	case 3:
		return StringWith(WithNumber|WithLowerLetter|WithUpperLetter|WithSpecial, c.Length)
	default:
		return StringWith(WithNumber, c.Length)
	}
}

func (c *Args) Int() int {
	minv := stringx.ParseInt(c.Min, 1)
	maxv := stringx.ParseInt(c.Max, 999)
	return IntRange(minv, maxv)
}

func (c *Args) Float() float64 {
	minv := stringx.ParseFloat(c.Min, 1)
	maxv := stringx.ParseFloat(c.Max, 999)
	prec := intx.IfZero(c.Prec, 6)
	return Float64Range(minv, maxv, prec)
}

func (c *Args) Time() time.Time {
	end := stringx.ParseTime(c.Max, time.Now())
	start := stringx.ParseTime(c.Min, end.Add(time.Hour*-24*30))
	return TimeRange(start, end)
}

func (c *Args) TimeFmt(format ...string) string {
	end := stringx.ParseTime(c.Max, time.Now())
	start := stringx.ParseTime(c.Min, end.Add(time.Hour*-24*30))
	var layout = timex.TimeFmt
	if len(format) > 0 {
		layout = format[0]
	} else if c.Format != "" {
		layout = c.Format
	}
	return TimeRange(start, start).Format(layout)
}

func (c *Args) Enum() string {
	return Enum(c.Enums)
}

func (c *Args) String() string {
	return String(c.Length)
}
