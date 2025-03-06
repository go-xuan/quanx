package maskx

import "strings"

type Desensitize uint

const mask = "*"

const (
	Phone Desensitize = iota
	Name
	IdCard
	Email
	Default // 全打码
)

func (d Desensitize) Desensitize(s string) string {
	if s != "" {
		switch {
		case d == Phone && len(s) == 11:
			s = s[:3] + strings.Repeat(mask, 4) + s[7:]
		case d == IdCard && len(s) == 18:
			s = s[:4] + strings.Repeat(mask, 11) + s[15:]
		case d == Name:
			split := strings.Split(s, "")
			s = split[0] + strings.Repeat(mask, len(split)-1)
		case d == Email:
			i := strings.Index(s, "@")
			s = strings.Repeat(mask, i) + s[i:]
		default:
			s = strings.Repeat(mask, len(s))
		}
	}
	return s
}
