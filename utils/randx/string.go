package randx

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/go-xuan/quanx/types/stringx"
)

// 随机字符串
func String(length ...int) string {
	var l = IntRange(5, 10)
	if len(length) > 0 && length[0] > 0 {
		l = length[0]
	}
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = SelectByte(allChar)
	}
	return string(bytes)
}

// 根据with生成字符串，可用于生成不同强度的密码
func StringWith(with int, length ...int) string {
	var l = IntRange(8, 16)
	if len(length) > 0 && length[0] > 0 {
		l = length[0]
	}
	bytes := make([]byte, l)
	var temp = numbers
	if with&WithLowerLetter > 0 {
		temp += lowerLetters
	}
	if with&WithUpperLetter > 0 {
		temp += upperLetters
	}
	if with&WithSpecial > 0 {
		temp += special
	}
	for i := 0; i < l; i++ {
		bytes[i] = SelectByte(temp)
	}
	for i := range bytes {
		j := NewRand().Intn(i + 1)
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return string(bytes)
}

// 随机长度数字码
func NumberCode(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		y := IntRange(0, len(numbers)-1)
		bytes[i] = numbers[y]
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
	sb.WriteString(Split(surname, ","))
	sb.WriteString(Split(numberCn, ","))
	sb.WriteString(Split(shengXiao, ","))
	return sb.String()
}

// 随机手机号
func Phone() string {
	bytes := make([]byte, 11)
	bytes[0] = '1'
	bytes[1] = phonePrefix[IntRange(0, len(phonePrefix)-1)]
	for i := 2; i < 11; i++ {
		y := IntRange(0, len(numbers)-1)
		bytes[i] = numbers[y]
	}
	return string(bytes)
}

// 随机身份证,湖北省内
func IdCard() string {
	sb := strings.Builder{}
	sb.WriteString(Split(hubeiProvinceCode, ","))
	sb.WriteString(time.Unix(Int64Range(1, time.Now().Unix()), 0).Format("20060102"))
	sb.WriteString(NumberCode(3))
	last := stringx.GetIdCardLastCode(sb.String())
	sb.WriteByte(last)
	return sb.String()
}

// 随机车牌号
func PlateNo() string {
	sb := strings.Builder{}
	sb.WriteString(Split(provinceSimple, ","))
	sb.WriteByte(SelectByte(upperLetters))
	for i := 0; i < 5; i++ {
		if Bool() {
			z := IntRange(0, len(upperLetters)-1)
			sb.WriteString(string(upperLetters[z]))
		} else {
			z := IntRange(0, len(numbers)-1)
			sb.WriteString(string(numbers[z]))
		}
	}
	return sb.String()
}

// 随机邮箱号
func Email() string {
	sb := strings.Builder{}
	x, y := IntRange(5, 10), IntRange(2, 5)
	for i := 0; i < x; i++ {
		sb.WriteByte(SelectByte(lowerChar))
	}
	sb.WriteString(`@`)
	for i := 0; i < y; i++ {
		sb.WriteByte(SelectByte(lowerLetters))
	}
	sb.WriteString(`.com`)
	return sb.String()
}

func IP() string {
	return fmt.Sprintf("%d.%d.%d.%d", IntRange(1, 255), IntRange(0, 255), IntRange(0, 255), IntRange(0, 255))
}

func Province() string {
	list := strings.Split(provinceName, ",")
	return list[IntRange(0, len(list)-1)]
}

func City() string {
	list := strings.Split(hubeiCityName, ",")
	return list[IntRange(0, len(list)-1)]
}
