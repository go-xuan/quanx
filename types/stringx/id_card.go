package stringx

import (
	"strconv"
	"time"
)

var weight = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
var IdCardLastCodes = [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
var provinceCodes = []string{
	"11", "12", "13", "14", "15",
	"21", "22", "23",
	"31", "32", "33", "34", "35", "36", "37",
	"41", "42", "43", "44", "45", "46",
	"50", "51", "52", "53", "54",
	"61", "62", "63", "64", "65",
	"71", "81", "91",
}

// GetIdCardLastCode 获生成身份证最后一位校验码
// 十七位计算公式：IdCardLastCodes[Sum(x*y)%11]
// x:表示第i位置上的身份证号码数字值
// y:表示第i位置上的加权因子
func GetIdCardLastCode(idCard string) byte {
	var array = make([]int, 17)
	for i := 0; i < 17; i++ {
		array[i], _ = strconv.Atoi(string(idCard[i]))
	}
	var sum int
	for i := 0; i < 17; i++ {
		sum += array[i] * weight[i]
	}
	return IdCardLastCodes[sum%11]
}

// GetIdCardGender 根据身份证获取性别
func GetIdCardGender(idCard string) string {
	if len(idCard) == 18 {
		sexSymbol, _ := strconv.Atoi(idCard[16:17])
		if sexSymbol%2 == 0 {
			return "女"
		} else {
			return "男"
		}
	}
	return "未知"
}

// GetIdCardAge 根据身份证获取年龄
func GetIdCardAge(idCard string) (age int) {
	year, _ := strconv.Atoi(idCard[6:10])
	month, _ := strconv.Atoi(idCard[10:12])
	day, _ := strconv.Atoi(idCard[12:14])
	birthday := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	age = time.Now().Year() - year
	if birthday.AddDate(age, 0, 0).After(time.Now()) {
		age--
	}
	return
}

// CheckIdCard 校验身份证
func CheckIdCard(idCard string) bool {
	length := len(idCard)
	if length != 15 && length != 18 {
		return false
	}
	return checkProvince(idCard) && checkBirthday(idCard) && checkLastCode(idCard)
}

// 校验出生年月日
func checkBirthday(idCard string) bool {
	year, _ := strconv.Atoi(idCard[6:10])
	month, _ := strconv.Atoi(idCard[10:12])
	day, _ := strconv.Atoi(idCard[12:14])
	if year < 1900 || month <= 0 || month > 12 || day <= 0 || day > 31 {
		return false
	}
	curYear, curMonth, curDay := time.Now().Date()
	if year == curYear {
		if month > int(curMonth) {
			return false
		} else if month == int(curMonth) && day > curDay {
			return false
		}
	}
	if month == 2 {
		if isLeapYear(year) && day > 29 {
			return false
		} else if day > 28 {
			return false
		}
	} else if 4 == month || 6 == month || 9 == month || 11 == month {
		if day > 30 {
			return false
		}
	}
	return true
}

// 校验省份
func checkProvince(idCard string) bool {
	provinceCode := make([]byte, 0)
	provinceCode = append(provinceCode, idCard[:2]...)
	provinceStr := string(provinceCode)
	for i := range provinceCodes {
		if provinceStr == provinceCodes[i] {
			return true
		}
	}
	return false
}

// 校验最后一位验证码是否正确
func checkLastCode(idCard string) bool {
	return GetIdCardLastCode(idCard) == idCard[17]
}

// 是否闰年
func isLeapYear(year int) bool {
	if year <= 0 {
		return false
	}
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return true
	}
	return false
}
