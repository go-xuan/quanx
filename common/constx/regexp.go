package constx

import "regexp"

const (
	ChineseRegex  = `^[\u4e00-\u9fa5]+$`                                 // 中文
	EnglishRegex  = `^[a-zA-Z]+$`                                        // 英文
	IntegerRegex  = `^[0-9]+$`                                           // 数字
	FloatRegex    = `^[0-9]+.[0-9]+$`                                    // 数字
	PhoneRegex    = `^1[3-9]\d{9}$`                                      // 手机号格式
	EmailRegex    = `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$` // 邮箱地址格式
	PasswordRegex = `^[a-zA-Z0-9_-]{6,18}$`                              // 密码格式
	DatetimeRegex = `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`              // 日期时间
	DateRegex     = `^\d{4}-\d{2}-\d{2}$`                                // 日期格式
	TimeRegex     = `^\d{2}:\d{2}:\d{2}$`                                // 时间格式
)

// 正则判断
func CheckRegexp(target, pattern string) bool {
	matched, _ := regexp.MatchString(pattern, target)
	return matched
}

func CheckChinese(s string) bool {
	return CheckRegexp(ChineseRegex, s)
}

func CheckEnglish(s string) bool {
	return CheckRegexp(EnglishRegex, s)
}

func CheckInteger(s string) bool {
	return CheckRegexp(IntegerRegex, s)
}

func CheckFloat(s string) bool {
	return CheckRegexp(FloatRegex, s)
}

func CheckPhone(s string) bool {
	return CheckRegexp(PhoneRegex, s)
}

func CheckEmail(s string) bool {
	return CheckRegexp(EmailRegex, s)
}

func CheckPassword(s string) bool {
	return CheckRegexp(PasswordRegex, s)
}
func CheckDatetime(s string) bool {
	return CheckRegexp(DatetimeRegex, s)
}

func CheckDate(s string) bool {
	return CheckRegexp(DateRegex, s)
}

func CheckTime(s string) bool {
	return CheckRegexp(TimeRegex, s)
}
