package randx

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-xuan/quanx/utilx/stringx"
	"github.com/google/uuid"
)

// 随机字符串
func String(size ...int) string {
	length := IntRange(5, 10)
	if len(size) > 0 {
		length = size[0]
	}
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		y := IntRange(0, len(AllChar)-1)
		bytes[i] = AllChar[y]
	}
	return string(bytes)
}

// 随机枚举
func Enum(enum string) string {
	list := strings.Split(enum, ",")
	return list[IntRange(0, len(list)-1)]
}

// 单选
func Radio(options []string) string {
	return options[IntRange(0, len(options)-1)]
}

// 随机长度数字码
func NumberCode(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		y := IntRange(0, len(Numbers)-1)
		bytes[i] = Numbers[y]
	}
	return string(bytes)
}

// 随机uuid
func UUID() string {
	return uuid.NewString()
}

// 随机长度数字码
func RandFrom(in string, sep string) string {
	split := strings.Split(in, sep)
	return split[IntRange(0, len(split)-1)]
}

// 随机姓名
func Name() string {
	sep := ","
	sb := strings.Builder{}
	sb.WriteString(RandFrom(XingShi, sep))
	sb.WriteString(RandFrom(NumberCn, sep))
	sb.WriteString(RandFrom(ShengXiao, sep))
	return sb.String()
}

// 随机手机号
func Phone() string {
	bytes := make([]byte, 11)
	bytes[0] = '1'
	x := IntRange(0, len(PhonePrefix)-1)
	bytes[1] = PhonePrefix[x]
	for i := 2; i < 11; i++ {
		y := IntRange(0, len(Numbers)-1)
		bytes[i] = Numbers[y]
	}
	return string(bytes)
}

// 随机身份证,湖北省内
func IdCard() string {
	sb := strings.Builder{}
	now := time.Now().Unix()
	diff := Int64Range(1, now)
	birthday := time.Unix(diff, 0).Format("20060102")
	provinceCodes := strings.Split(HubeiProvincecode, ",")
	x := IntRange(0, len(provinceCodes)-1)
	sb.WriteString(provinceCodes[x])
	sb.WriteString(birthday)
	sb.WriteString(NumberCode(3))
	last := stringx.GetIdCardLastCode(sb.String())
	sb.WriteByte(last)
	return sb.String()
}

// 随机车牌号
func PlateNo() string {
	sb := strings.Builder{}
	provinces := strings.Split(ProvinceSimple, ",")
	x := IntRange(0, len(provinces)-1)
	y := IntRange(0, len(Uppers)-20)
	sb.WriteString(provinces[x])
	sb.WriteString(string(Uppers[y]))
	for i := 0; i < 5; i++ {
		if Bool() {
			z := IntRange(0, len(Uppers)-1)
			sb.WriteString(string(Uppers[z]))
		} else {
			z := IntRange(0, len(Numbers)-1)
			sb.WriteString(string(Numbers[z]))
		}
	}
	return sb.String()
}

// 随机邮箱号
func Email() string {
	sb := strings.Builder{}
	len1 := IntRange(5, 10)
	for i := 0; i < len1; i++ {
		x := IntRange(0, len(LowerChar)-1)
		sb.WriteString(string(LowerChar[x]))
	}
	sb.WriteString(`@`)
	len2 := IntRange(2, 5)
	for i := 0; i < len2; i++ {
		x := IntRange(0, len(Lowers)-1)
		sb.WriteString(string(Lowers[x]))
	}
	sb.WriteString(`.com`)
	return sb.String()
}

func IP() string {
	ip := fmt.Sprintf("%d.%d.%d.%d",
		IntRange(1, 255),
		IntRange(0, 255),
		IntRange(0, 255),
		IntRange(0, 255))
	return ip
}

func Province() string {
	list := strings.Split(ProvinceName, ",")
	return list[IntRange(0, len(list)-1)]
}

func City() string {
	list := strings.Split(HubeiCityName, ",")
	return list[IntRange(0, len(list)-1)]
}

// 生成随机密码
func Password(length int, contains ...string) string {
	if length < 8 {
		length = 8
	}
	bytes := make([]byte, length)
	var temp string
	if len(contains) > 0 {
		for _, contain := range contains {
			switch contain {
			case Upper:
				temp += Uppers
			case Lower:
				temp += Lowers
			case Symbol:
				temp += SPECIAL
			case Number:
				temp += Numbers
			default:
			}
		}
	} else {
		temp = AllChar
	}

	for i := 0; i < length; i++ {
		x := IntRange(0, len(temp)-1)
		bytes[i] = temp[x]
	}
	for index := range bytes {
		j := NewRand().Intn(index + 1)
		bytes[index], bytes[j] = bytes[j], bytes[index]
	}
	return string(bytes)
}
