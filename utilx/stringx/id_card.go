package stringx

import (
	"strconv"
	"time"
)

var Weight = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
var ValidCodes = [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
var ProvinceCodes = []string{
	"11", "12", "13", "14", "15",
	"21", "22", "23",
	"31", "32", "33", "34", "35", "36", "37",
	"41", "42", "43", "44", "45", "46",
	"50", "51", "52", "53", "54",
	"61", "62", "63", "64", "65",
	"71", "81", "91",
}

// 校验身份证
func CheckIdCard(idCard string) bool {
	length := len(idCard)
	if length != 15 && length != 18 {
		return false
	}
	// 校验省份编码
	if !CheckProvince(idCard) {
		return false
	}
	// 验证出生日期
	year, _ := strconv.Atoi(idCard[6:10])
	month, _ := strconv.Atoi(idCard[10:12])
	day, _ := strconv.Atoi(idCard[12:14])
	if !checkBirthday(year, month, day) {
		return false
	}
	// 校验最后一位校验码
	if !CheckIdCard18th(idCard) {
		return false
	}
	return true
}

// 校验出生年月日
func checkBirthday(year, month, day int) bool {
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
	if 2 == month {
		if IsLeapYear(year) && day > 29 {
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
func CheckProvince(idCard string) bool {
	provinceCode := make([]byte, 0)
	provinceCode = append(provinceCode, idCard[:2]...)
	provinceStr := string(provinceCode)
	for i := range ProvinceCodes {
		if provinceStr == ProvinceCodes[i] {
			return true
		}
	}
	return false
}

// 校验第18位验证码是否正确
func CheckIdCard18th(idCard18 string) bool {
	return GetIdCard18th(idCard18) == idCard18[17]
}

// 十七位计算公式 N=Sum(x*y)
// x:表示第i位置上的身份证号码数字值
// y:表示第i位置上的加权因子
func GetIdCard18th(idCard string) byte {
	array := make([]int, 17)
	for i := 0; i < 17; i++ {
		array[i], _ = strconv.Atoi(string(idCard[i]))
	}
	var res int
	for i := 0; i < 17; i++ {
		res += array[i] * Weight[i]
	}
	return ValidCodes[res%11]
}

// 是否闰年
func IsLeapYear(year int) bool {
	if year <= 0 {
		return false
	}
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return true
	}
	return false
}

// 根据身份证获取性别
func GetGender(idNo string) (gender string) {
	if len(idNo) == 18 {
		sexSymbol, _ := strconv.Atoi(idNo[16:17])
		if sexSymbol%2 == 0 {
			gender = "女"
		} else {
			gender = "男"
		}
	}
	return
}

// 根据身份证获取年龄
func GetAge(idNo string) (age int) {
	year, _ := strconv.Atoi(idNo[6:10])
	month, _ := strconv.Atoi(idNo[10:12])
	day, _ := strconv.Atoi(idNo[12:14])
	birthday := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	age = time.Now().Year() - year
	if birthday.AddDate(age, 0, 0).After(time.Now()) {
		age--
	}
	return
}
