package stringx

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
func RegMatch(target, pattern string) bool {
	matched, _ := regexp.MatchString(pattern, target)
	return matched
}

func IsChinese(s string) bool {
	return RegMatch(ChineseRegex, s)
}

func IsEnglish(s string) bool {
	return RegMatch(EnglishRegex, s)
}

func IsInteger(s string) bool {
	return RegMatch(IntegerRegex, s)
}

func IsFloat(s string) bool {
	return RegMatch(FloatRegex, s)
}

func IsPhone(s string) bool {
	return RegMatch(PhoneRegex, s)
}

func IsEmail(s string) bool {
	return RegMatch(EmailRegex, s)
}

func IsPassword(s string) bool {
	return RegMatch(PasswordRegex, s)
}

func IsDatetime(s string) bool {
	return RegMatch(DatetimeRegex, s)
}

func IsDate(s string) bool {
	return RegMatch(DateRegex, s)
}

func IsTime(s string) bool {
	return RegMatch(TimeRegex, s)
}
