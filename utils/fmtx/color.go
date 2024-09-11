package fmtx

import (
	"fmt"
)

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	Grey
)

type Color uint8

func (c Color) String(s string) string {
	if Black <= c && c <= Grey {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
	}
	return s
}

func (c Color) Bytes(s string) []byte {
	if Black <= c && c <= Grey {
		return []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s))
	}
	return []byte(s)
}

func (c Color) Println(s string) {
	fmt.Println(c.String(s))
}

// Printf 打印
func (c Color) Printf(s string, v ...any) {
	fmt.Println(c.String(fmt.Sprintf(s, v...)))
}

// XPrintf 打印时仅对可变参数生效
func (c Color) XPrintf(s string, v ...any) {
	if len(v) > 0 {
		for i, a := range v {
			v[i] = c.String(fmt.Sprint(a))
		}
		fmt.Printf(s, v...)
		fmt.Println()
	} else {
		fmt.Println(s)
	}
}

// Sprintf 格式化
func (c Color) Sprintf(s string, v ...any) string {
	return fmt.Sprintf(c.String(s), v...)
}

// XSPrintf 格式化仅对可变参数生效
func (c Color) XSPrintf(s string, v ...any) string {
	if len(v) > 0 {
		for i, a := range v {
			v[i] = c.String(fmt.Sprint(a))
		}
		return fmt.Sprintf(s, v...)
	} else {
		return fmt.Sprint(s)
	}
}
