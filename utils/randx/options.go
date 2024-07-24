package randx

import (
	"strings"
	"time"

	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/types/timex"
)

const (
	randTypeInt       = "int"        // 数字
	randTypeFloat     = "float"      // 浮点数
	randTypeSequence  = "sequence"   // 数字编号
	randTypeTime      = "time"       // 时间
	randTypeDate      = "date"       // 日期
	randTypeUUID      = "uuid"       // uuid
	randTypePhone     = "phone"      // 手机号
	randTypeName      = "name"       // 姓名
	randTypeIdCard    = "id_card"    // 身份证
	randTypePlateNo   = "plate_no"   // 车牌号
	randTypeEmail     = "email"      // 邮箱
	randTypeIP        = "ip"         // ip地址
	randTypeProvince  = "province"   // 省
	randTypeCity      = "city"       // 市
	randTypePassword  = "password"   // 密码
	randTypeEnum      = "enum"       // 枚举
	randTypeOtherEnum = "other_enum" // 其他枚举
)

// Options 随机生成
type Options struct {
	Type    string   // 数据类型
	Param   *Param   // 约束条件参数
	Default string   // 默认值
	Offset  int      // 偏移量
	Enums   []string // 枚举
}

// RandData 生成随机数
func (o *Options) RandData() any {
	switch o.Type {
	case randTypeInt:
		return o.RandInt()
	case randTypeFloat:
		return o.RandFloat()
	default:
		return o.RandString()
	}
}

// RandDataString 生成随机数
func (o *Options) RandDataString() string {
	switch o.Type {
	case randTypeInt:
		return stringx.ParseInt(o.RandInt())
	case randTypeFloat:
		return stringx.ParseFloat64(o.RandFloat())
	default:
		return o.RandString()
	}
}

// RandInt 生成随机数字
func (o *Options) RandInt() int {
	if o.Default == "" {
		return o.Param.Int()
	} else {
		return stringx.ToInt(o.Default)
	}
}

// RandFloat 生成随机浮点数
func (o *Options) RandFloat() float64 {
	if o.Default == "" {
		return o.Param.Float()
	} else {
		return stringx.ToFloat64(o.Default)
	}
}

// RandString 生成随机字符串
func (o *Options) RandString() (result string) {
	if param, def := o.Param, o.Default; param != nil && def == "" {
		if o.Type == randTypeSequence {
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
	case randTypePhone:
		return Phone()
	case randTypeName:
		return Name()
	case randTypeIdCard:
		return IdCard()
	case randTypePlateNo:
		return PlateNo()
	case randTypeEmail:
		return Email()
	case randTypeUUID:
		return UUID()
	case randTypeIP:
		return IP()
	case randTypeProvince:
		return Province()
	case randTypeCity:
		return City()
	case randTypePassword:
		return o.Param.Password()
	case randTypeDate:
		return o.Param.TimeFmt(timex.DateFmt)
	case randTypeTime:
		return o.Param.TimeFmt()
	case randTypeEnum:
		return Enum(o.Param.Enums)
	case randTypeOtherEnum:
		return Enum(append(o.Enums, o.Param.Enums...))
	default:
		return String(o.Param.Length)
	}
}

// Param 随机数生成参数
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
		return StringWith(WithNumber|WithLowerLetter|WithUpperLetter, c.Length)
	case 3:
		return StringWith(WithNumber|WithLowerLetter|WithUpperLetter|WithSpecial, c.Length)
	default:
		return StringWith(WithNumber, c.Length)
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
