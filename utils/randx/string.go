package randx

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/go-xuan/quanx/types/stringx"
)

// 随机字符串
func String(size ...int) string {
	length := IntRange(5, 10)
	if len(size) > 0 && size[0] > 0 {
		length = size[0]
	}
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		y := IntRange(0, len(allChar)-1)
		bytes[i] = allChar[y]
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
	sb.WriteString(Split(xingShi, ","))
	sb.WriteString(Split(numberCn, ","))
	sb.WriteString(Split(shengXiao, ","))
	return sb.String()
}

// 随机手机号
func Phone() string {
	bytes := make([]byte, 11)
	bytes[0] = '1'
	x := IntRange(0, len(phonePrefix)-1)
	bytes[1] = phonePrefix[x]
	for i := 2; i < 11; i++ {
		y := IntRange(0, len(numbers)-1)
		bytes[i] = numbers[y]
	}
	return string(bytes)
}

// 随机身份证,湖北省内
func IdCard() string {
	sb := strings.Builder{}
	now := time.Now().Unix()
	diff := Int64Range(1, now)
	birthday := time.Unix(diff, 0).Format("20060102")
	sb.WriteString(Split(hubeiProvinceCode, ","))
	sb.WriteString(birthday)
	sb.WriteString(NumberCode(3))
	last := stringx.GetIdCardLastCode(sb.String())
	sb.WriteByte(last)
	return sb.String()
}

// 随机车牌号
func PlateNo() string {
	sb := strings.Builder{}
	y := IntRange(0, len(upperLetters)-20)
	sb.WriteString(Split(provinceSimple, ","))
	sb.WriteString(string(upperLetters[y]))
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
	len1 := IntRange(5, 10)
	for i := 0; i < len1; i++ {
		x := IntRange(0, len(lowerChar)-1)
		sb.WriteString(string(lowerChar[x]))
	}
	sb.WriteString(`@`)
	len2 := IntRange(2, 5)
	for i := 0; i < len2; i++ {
		x := IntRange(0, len(lowerLetters)-1)
		sb.WriteString(string(lowerLetters[x]))
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
	list := strings.Split(provinceName, ",")
	return list[IntRange(0, len(list)-1)]
}

func City() string {
	list := strings.Split(hubeiCityName, ",")
	return list[IntRange(0, len(list)-1)]
}

const (
	numberPwd      = 1 << 0 // 数字
	lowerLetterPwd = 1 << 1 // 小写字母
	upperLetterPwd = 1 << 2 // 大写字母
	specialPwd     = 1 << 3 // 特殊符号
	SimplePwd      = numberPwd
	MediumPwd      = numberPwd | lowerLetterPwd | upperLetterPwd
	HardPwd        = numberPwd | lowerLetterPwd | upperLetterPwd | specialPwd
)

// 生成随机密码
func Password(length int, level int) string {
	if length < 8 {
		length = 8
	}
	bytes := make([]byte, length)
	var temp = numbers
	if level&lowerLetterPwd > 0 {
		temp += lowerLetters
	}
	if level&upperLetterPwd > 0 {
		temp += upperLetters
	}
	if level&specialPwd > 0 {
		temp += special
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
