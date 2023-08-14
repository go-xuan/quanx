package randx

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/quanxiaoxuan/quanx/utils/stringx"
)

// 随机字符串
func String() string {
	length := IntRange(5, 10)
	return CharString(length)
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

// 随机长度字符串
func CharString(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		y := IntRange(0, len(CharAll)-1)
		bytes[i] = CharAll[y]
	}
	return string(bytes)
}

// 随机长度数字码
func NumberString(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		y := IntRange(0, len(NUMBER)-1)
		bytes[i] = NUMBER[y]
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
	familyNames := strings.Split(FamilyNameCn, ",")
	numbers := strings.Split(NumberCn, ",")
	animals := strings.Split(Animal, ",")
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
	x := IntRange(0, len(PhonePrefix)-1)
	bytes[1] = PhonePrefix[x]
	for i := 2; i < 11; i++ {
		y := IntRange(0, len(NUMBER)-1)
		bytes[i] = NUMBER[y]
	}
	return string(bytes)
}

// 随机身份证,湖北省内
func IdCard() string {
	sb := strings.Builder{}
	now := time.Now().Unix()
	diff := Int64Range(1, now)
	birthday := time.Unix(diff, 0).Format("20060102")
	postCodes := strings.Split(HubeiPostcode, ",")
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
	provinces := strings.Split(ProvinceSimple, ",")
	x := IntRange(0, len(provinces)-1)
	y := IntRange(0, len(LetterUpper)-20)
	sb.WriteString(provinces[x])
	sb.WriteString(string(LetterUpper[y]))
	for i := 0; i < 5; i++ {
		if Bool() {
			z := IntRange(0, len(LetterUpper)-1)
			sb.WriteString(string(LetterUpper[z]))
		} else {
			z := IntRange(0, len(NUMBER)-1)
			sb.WriteString(string(NUMBER[z]))
		}
	}
	return sb.String()
}

// 随机邮箱号
func Email() string {
	sb := strings.Builder{}
	len1 := IntRange(5, 10)
	for i := 0; i < len1; i++ {
		x := IntRange(0, len(CharLower)-1)
		sb.WriteString(string(CharLower[x]))
	}
	sb.WriteString(`@`)
	len2 := IntRange(2, 5)
	for i := 0; i < len2; i++ {
		x := IntRange(0, len(LetterLower)-1)
		sb.WriteString(string(LetterLower[x]))
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
func Password(length int, lower bool, upper bool, numeric bool, special bool) string {
	if length < 8 {
		length = 8
	}
	bytes := make([]byte, length)
	var temp string
	if lower {
		temp += LetterLower
	}
	if upper {
		temp += LetterUpper
	}
	if numeric {
		temp += NUMBER
	}
	if special {
		temp += SPECIAL
	}
	if temp == "" {
		temp = CharLower
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
