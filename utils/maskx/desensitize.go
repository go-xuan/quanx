package maskx

import "strings"

type Desensitize uint

const mask = "*"

const (
	Phone   Desensitize = iota // 手机号，中间4位打码，18029066575 ==> 180****6575
	IdCard                     // 身份证，中间11位打码，例如：420300198308109549 ==> 4203***********549
	Email                      // 邮箱，保留后缀，例如：zhangsan@jwzg.com ==> ********@jwzg.com
	Name                       // 姓名，保留姓氏，例如：张三三 ==> 张**
	Default                    // 全打码
)

// Desensitize 脱敏
func (d Desensitize) Desensitize(text string) string {
	switch {
	case d == Phone && len(text) == 11:
		return text[:3] + strings.Repeat(mask, 4) + text[7:]
	case d == IdCard && len(text) == 18:
		return text[:4] + strings.Repeat(mask, 11) + text[15:]
	case d == Email:
		i := strings.Index(text, "@")
		return strings.Repeat(mask, i) + text[i:]
	case d == Name:
		split := strings.Split(text, "")
		return split[0] + strings.Repeat(mask, len(split)-1)
	default:
		return strings.Repeat(mask, len(text))
	}
}
