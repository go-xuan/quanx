package randx

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-xuan/quanx/types/intx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/types/timex"
)

func NewParam(data map[string]string) *Param {
	return &Param{
		Min:    data["min"],
		Max:    data["max"],
		Prefix: data["prefix"],
		Suffix: data["suffix"],
		Upper:  data["upper"] == "true",
		Lower:  data["lower"] == "true",
		Old:    data["old"],
		New:    data["new"],
		Format: data["format"],
		Length: stringx.ParseInt(data["length"]),
		Prec:   stringx.ParseInt(data["prec"]),
		Level:  stringx.ParseInt(data["level"]),
		Enums:  strings.Split(data["enums"], ","),
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

// 字符串
func (p *Param) String(t string) string {
	switch t {
	case typeInt:
		return strconv.Itoa(p.Int())
	case typeFloat:
		return strconv.FormatFloat(p.Float(), 'f', -1, 64)
	case typeTime:
		return p.TimeFmt()
	case typeDate:
		return p.TimeFmt(timex.DateFmt)
	case typeUUID:
		return UUID()
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
	case typeIP:
		return IP()
	case typeProvince:
		return Province()
	case typeCity:
		return City()
	case typePassword:
		return p.Password()
	case typeEnum:
		return StringFrom(p.Enums...)
	default:
		return String(p.Length)
	}
}

// Int 整数
func (p *Param) Int() int {
	min := stringx.ParseInt(p.Min, 1)
	max := stringx.ParseInt(p.Max, 999)
	return IntRange(min, max)
}

// Float 浮点数
func (p *Param) Float() float64 {
	min := stringx.ParseFloat(p.Min, 1)
	max := stringx.ParseFloat(p.Max, 999)
	prec := intx.IfZero(p.Prec, 6)
	return Float64Range(min, max, prec)
}

// Time 时间
func (p *Param) Time() time.Time {
	start := stringx.ParseTime(p.Min, time.Time{})
	end := stringx.ParseTime(p.Max, time.Now())
	return TimeRange(start, end)
}

// TimeFmt 时间字符串
func (p *Param) TimeFmt(format ...string) string {
	start := stringx.ParseTime(p.Min, time.Time{})
	end := stringx.ParseTime(p.Max, time.Now())
	var layout = timex.TimeFmt
	if len(format) > 0 {
		layout = format[0]
	} else if p.Format != "" {
		layout = p.Format
	}
	return TimeRange(start, end).Format(layout)
}

// Enum 枚举值
func (p *Param) Enum() string {
	return StringFrom(p.Enums...)
}

// Sequence 序列值
func (p *Param) Sequence(offset int) int {
	return stringx.ParseInt(p.Min, 1) + offset
}

// Password 密码
func (p *Param) Password() string {
	switch p.Level {
	case 2:
		return StringUse(UseNumber|UseLowerLetter|UseUpperLetter, p.Length)
	case 3:
		return StringUse(UseNumber|UseLowerLetter|UseUpperLetter|UseSpecialSymbols, p.Length)
	default:
		return StringUse(UseNumber, p.Length)
	}
}

// AddEnum 添加枚举项
func (p *Param) AddEnum(enum ...string) {
	p.Enums = append(p.Enums, enum...)
}

// Modify 修饰
func (p *Param) Modify(value string) string {
	if p.Old != "" {
		value = strings.ReplaceAll(value, p.Old, p.New)
	}
	if p.Prefix != "" {
		value = p.Prefix + value
	}
	if p.Suffix != "" {
		value = value + p.Suffix
	}
	if p.Upper {
		value = strings.ToUpper(value)
	}
	if p.Lower {
		value = strings.ToLower(value)
	}
	return value
}
