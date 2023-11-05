package randx

import (
	"fmt"
	"github.com/go-xuan/quanx/common/constx"
	"strings"
	"time"

	"github.com/go-xuan/quanx/utils/stringx"
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
		y := IntRange(0, len(constx.CharAll)-1)
		bytes[i] = constx.CharAll[y]
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
func NumberString(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		y := IntRange(0, len(constx.NUMBER)-1)
		bytes[i] = constx.NUMBER[y]
	}
	return string(bytes)
}

// 随机uuid
func UUID() string {
	return uuid.NewString()
}

// 随机姓名
func Name() string {
	sb := strings.Builder{}
	familyNames := strings.Split(constx.FamilyNameCn, ",")
	numbers := strings.Split(constx.NumberCn, ",")
	animals := strings.Split(constx.Animal, ",")
	x := IntRange(0, len(familyNames)-1)
	y := IntRange(0, len(numbers)-1)
	z := IntRange(0, len(animals)-1)
	sb.WriteString(familyNames[x])
	sb.WriteString(numbers[y])
	sb.WriteString(animals[z])
	return sb.String()
}

// 随机手机号
func Phone() string {
	bytes := make([]byte, 11)
	bytes[0] = '1'
	x := IntRange(0, len(constx.PhonePrefix)-1)
	bytes[1] = constx.PhonePrefix[x]
	for i := 2; i < 11; i++ {
		y := IntRange(0, len(constx.NUMBER)-1)
		bytes[i] = constx.NUMBER[y]
	}
	return string(bytes)
}

// 随机身份证,湖北省内
func IdCard() string {
	sb := strings.Builder{}
	now := time.Now().Unix()
	diff := Int64Range(1, now)
	birthday := time.Unix(diff, 0).Format("20060102")
	postCodes := strings.Split(constx.HubeiPostcode, ",")
	x := IntRange(0, len(postCodes)-1)
	sb.WriteString(postCodes[x])
	sb.WriteString(birthday)
	sb.WriteString(NumberString(3))
	last := stringx.GetIdCard18th(sb.String())
	sb.WriteByte(last)
	return sb.String()
}

// 随机车牌号
func PlateNo() string {
	sb := strings.Builder{}
	provinces := strings.Split(constx.ProvinceSimple, ",")
	x := IntRange(0, len(provinces)-1)
	y := IntRange(0, len(constx.LetterUpper)-20)
	sb.WriteString(provinces[x])
	sb.WriteString(string(constx.LetterUpper[y]))
	for i := 0; i < 5; i++ {
		if Bool() {
			z := IntRange(0, len(constx.LetterUpper)-1)
			sb.WriteString(string(constx.LetterUpper[z]))
		} else {
			z := IntRange(0, len(constx.NUMBER)-1)
			sb.WriteString(string(constx.NUMBER[z]))
		}
	}
	return sb.String()
}

// 随机邮箱号
func Email() string {
	sb := strings.Builder{}
	len1 := IntRange(5, 10)
	for i := 0; i < len1; i++ {
		x := IntRange(0, len(constx.CharLower)-1)
		sb.WriteString(string(constx.CharLower[x]))
	}
	sb.WriteString(`@`)
	len2 := IntRange(2, 5)
	for i := 0; i < len2; i++ {
		x := IntRange(0, len(constx.LetterLower)-1)
		sb.WriteString(string(constx.LetterLower[x]))
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
	list := strings.Split(constx.ProvinceName, ",")
	return list[IntRange(0, len(list)-1)]
}

func City() string {
	list := strings.Split(constx.HubeiCityName, ",")
	return list[IntRange(0, len(list)-1)]
}

// 生成随机密码
func Password(length int, lower bool, upper bool, numeric bool, special bool) string {
	if length < 8 {
		length = 8
	}
	bytes := make([]byte, length)
	var temp string
	if lower {
		temp += constx.LetterLower
	}
	if upper {
		temp += constx.LetterUpper
	}
	if numeric {
		temp += constx.NUMBER
	}
	if special {
		temp += constx.SPECIAL
	}
	if temp == "" {
		temp = constx.CharLower
	}
	for i := 0; i < length; i++ {
		x := IntRange(0, len(temp)-1)
		bytes[i] = temp[x]
	}
	for index := range bytes {
		j := newRand.Intn(index + 1)
		bytes[index], bytes[j] = bytes[j], bytes[index]
	}
	return string(bytes)
}
