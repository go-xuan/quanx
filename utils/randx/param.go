package randx

import (
	"strings"
	"time"

	"github.com/go-xuan/quanx/types/intx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/types/timex"
)

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

func (p *Param) Int() int {
	minv := stringx.ParseInt(p.Min, 1)
	maxv := stringx.ParseInt(p.Max, 999)
	return IntRange(minv, maxv)
}

func (p *Param) Float() float64 {
	minv := stringx.ParseFloat(p.Min, 1)
	maxv := stringx.ParseFloat(p.Max, 999)
	prec := intx.IfZero(p.Prec, 6)
	return Float64Range(minv, maxv, prec)
}

func (p *Param) Time() time.Time {
	end := stringx.ParseTime(p.Max, time.Now())
	start := stringx.ParseTime(p.Min, time.Unix(0, 0))
	return TimeRange(start, end)
}

func (p *Param) TimeFmt(format ...string) string {
	end := stringx.ParseTime(p.Max, time.Now())
	start := stringx.ParseTime(p.Min, time.Unix(0, 0))
	var layout = timex.TimeFmt
	if len(format) > 0 {
		layout = format[0]
	} else if p.Format != "" {
		layout = p.Format
	}
	return TimeRange(start, end).Format(layout)
}

func (p *Param) Enum() string {
	return StringFrom(p.Enums...)
}

func (p *Param) String() string {
	return String(p.Length)
}
