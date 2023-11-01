package randx

import (
	"github.com/quanxiaoxuan/quanx/common/constx"
	"strconv"
	"strings"
	"time"

	"github.com/quanxiaoxuan/quanx/utils/stringx"
	"github.com/quanxiaoxuan/quanx/utils/timex"
)

// 随机数模型
type RandModel struct {
	Type       RandType `json:"type"`       // 数据类型
	Size       int      `json:"size"`       // 生成数量
	Default    string   `json:"default"`    // 默认值
	Constraint string   `json:"constraint"` // 约束条件限制,格式为"key=value&key=value&..."
	Index      int      `json:"index"`      // 下标
	Options    []string `json:"options"`    // 备选项
}

// 生成随机数
func (m *RandModel) GetRandData() (result any) {
	switch m.Type {
	case IntType:
		result = m.GetRandInt() // 数字
	case FloatType:
		result = m.GetRandFloat() // 浮点数
	default:
		result = m.GetRandString() // 字符穿
	}
	return
}

// 生成随机数字
func (m *RandModel) GetRandInt() (result int) {
	if m.Default == "" {
		result = IntRange(ParseConstraintOfInt(m.Constraint))
	} else {
		result, _ = strconv.Atoi(m.Default)
	}
	return
}

// 生成随机浮点数
func (m *RandModel) GetRandFloat() (result float64) {
	if m.Default == "" {
		result = Float64Range(ParseConstraintOfFloat(m.Constraint))
	} else {
		result, _ = strconv.ParseFloat(m.Default, 64)
	}
	return
}

// 生成随机字符串
func (m *RandModel) GetRandString() (result string) {
	if m.Default == "" {
		if m.Type == NoType {
			min, _ := ParseConstraintOfInt(m.Constraint)
			result = strconv.Itoa(min + m.Index)
		} else {
			result = m.randString()
		}
		kvMap := ParseConstraint(m.Constraint)
		// 补充前缀
		if strings.Contains(m.Constraint, constx.Prefix) {
			result = kvMap[constx.Prefix] + result
		}
		// 后缀
		if strings.Contains(m.Constraint, constx.Suffix) {
			result = result + kvMap[constx.Suffix]
		}
		// 字符替换
		if stringx.ContainsAny(m.Constraint, constx.Old, constx.New) {
			result = strings.ReplaceAll(result, kvMap[constx.Old], kvMap[constx.New])
		}
		// 大小写转换
		if strings.Contains(m.Constraint, constx.Upper) {
			if kvMap[constx.Upper] == "true" {
				result = strings.ToUpper(result)
			}
		}
		// 大小写转换
		if strings.Contains(m.Constraint, constx.Lower) {
			if kvMap[constx.Lower] == "true" {
				result = strings.ToLower(result)
			}
		}
	} else {
		result = m.Default
	}
	return
}

// 生成随机字符串
func (m *RandModel) randString() (result string) {
	constraint := m.Constraint
	switch m.Type {
	case PhoneType:
		result = Phone()
	case NameType:
		result = Name()
	case IdCardType:
		result = IdCard()
	case PlateNoType:
		result = PlateNo()
	case EmailType:
		result = Email()
	case UuidType:
		result = UUID()
	case IPType:
		result = IP()
	case ProvinceType:
		result = Province()
	case CityType:
		result = City()
	case PasswordType:
		result = Password(ParseConstraintOfPassword(constraint))
	case IntStringType:
		result = strconv.Itoa(IntRange(ParseConstraintOfInt(constraint)))
	case DateType:
		min, max, _ := ParseConstraintOfTime(constraint)
		result = DateRange(min, max)
	case TimeType:
		min, max, format := ParseConstraintOfTime(constraint)
		result = TimeRange(min, max).Format(format)
	case OptionType:
		result = Radio(ParseConstraintOfOptions(constraint))
	case DatabaseType:
		options := ParseConstraintOfOptions(constraint)
		options = append(options, m.Options...)
		result = Radio(options)
	case StringType:
		result = String(GetLength(constraint))
	default:
		result = String()
	}
	return
}

// 解析约束参数，转为map集合
func ParseConstraint(constraint string) map[string]string {
	if strings.Contains(constraint, constx.Equal) {
		kvMap := make(map[string]string)
		kvs := strings.Split(constraint, constx.Split)
		for _, kv := range kvs {
			k, v := stringx.SplitByFirst(kv, constx.Equal)
			kvMap[k] = v
		}
		return kvMap
	}
	return nil
}

// 获取长度约束
func GetLength(constraint string) (length int) {
	kvMap := ParseConstraint(constraint)
	if strings.Contains(constraint, constx.Length) {
		length, _ = strconv.Atoi(kvMap[constx.Length])
	}
	return 8
}

// 获取前后缀
func ParsePrefixAndSuffix(constraint string) (prefix, suffix string) {
	kvMap := ParseConstraint(constraint)
	if strings.Contains(constraint, constx.Prefix) {
		prefix = kvMap[constx.Prefix]
	}
	if strings.Contains(constraint, constx.Suffix) {
		suffix = kvMap[constx.Suffix]
	}
	return
}

// 解析数字约束条件
// 默认取值范围：1~9999
func ParseConstraintOfInt(constraint string) (min, max int) {
	min = 1
	max = 9999
	kvMap := ParseConstraint(constraint)
	if strings.Contains(constraint, constx.Min) {
		min, _ = strconv.Atoi(kvMap[constx.Min])
	}
	if strings.Contains(constraint, constx.Max) {
		max, _ = strconv.Atoi(kvMap[constx.Max])
	}
	return
}

// 解析浮点数生成范围以及精度
// 默认取值范围：1~9999
// 默认精度：保留6位小数
func ParseConstraintOfFloat(constraint string) (min, max float64, prec int) {
	min = 1
	max = 9999
	prec = 6
	kvMap := ParseConstraint(constraint)
	if strings.Contains(constraint, constx.Min) {
		min, _ = strconv.ParseFloat(kvMap[constx.Min], 64)
	}
	if strings.Contains(constraint, constx.Max) {
		max, _ = strconv.ParseFloat(kvMap[constx.Max], 64)
	}
	if strings.Contains(constraint, constx.Prec) {
		prec, _ = strconv.Atoi(kvMap[constx.Prec])
	}
	return
}

// 解析时间约束条件
// 默认取值范围：近30天
func ParseConstraintOfTime(constraint string) (min, max time.Time, format string) {
	now := time.Now()
	max = now
	min = now.Add(time.Hour * -24 * 30)
	format = constx.TimeFmt
	kvMap := ParseConstraint(constraint)
	if strings.Contains(constraint, constx.Min) {
		min = timex.ToTime(kvMap[constx.Min])
	}
	if strings.Contains(constraint, constx.Max) {
		max = timex.ToTime(kvMap[constx.Max])
	}
	if strings.Contains(constraint, constx.Format) {
		format = kvMap[constx.Format]
	}
	return
}

// 解析密码约束条件
func ParseConstraintOfPassword(constraint string) (length int, lower, upper, numeric, special bool) {
	if stringx.ContainsAny(constraint, constx.Length, constx.Lower, constx.Upper, constx.HasNumber, constx.HasSymbol) {
		kvMap := ParseConstraint(constraint)
		length, _ = strconv.Atoi(kvMap[constx.Length])
		lower = kvMap[constx.Lower] == "true"
		upper = kvMap[constx.Upper] == "true"
		numeric = kvMap[constx.HasNumber] == "true"
		special = kvMap[constx.HasSymbol] == "true"
		return
	}
	return 8, true, true, true, true
}

// 解析备选项约束条件
func ParseConstraintOfOptions(constraint string) (options []string) {
	if stringx.ContainsAny(constraint, constx.Options) {
		kvMap := ParseConstraint(constraint)
		options = strings.Split(kvMap[constx.Options], ",")
	}
	return
}

// 获取数据库取值sql
func ParseConstraintOfSql(constraint string) (sql string) {
	if stringx.ContainsAny(constraint, constx.Table, constx.Field) {
		sb := strings.Builder{}
		kvMap := ParseConstraint(constraint)
		sb.WriteString(`select distinct `)
		sb.WriteString(kvMap[constx.Field])
		sb.WriteString(` from `)
		sb.WriteString(kvMap[constx.Table])
		sql = sb.String()
		return
	}
	return
}
