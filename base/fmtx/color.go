package fmtx

import (
	"fmt"
)

const (
	Black   Color = iota + 30 // 黑色
	Red                       // 红色
	Green                     // 绿色
	Yellow                    // 黄色
	Blue                      // 蓝色
	Magenta                   // 洋红色
	Cyan                      // 青色
	Grey                      // 灰色
)

type Color uint8

// String 返回带有颜色格式的字符串
// 如果颜色值在有效范围内（ Black 到 Grey ），则添加 ANSI 颜色转义序列；否则返回原始字符串
func (c Color) String(s string) string {
	if Black <= c && c <= Grey {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
	}
	return s
}

// Bytes 返回带有颜色格式的字节切片
// 如果颜色值在有效范围内（ Black 到 Grey ），则添加 ANSI 颜色转义序列；否则返回原始字符串的字节切片
func (c Color) Bytes(s string) []byte {
	if Black <= c && c <= Grey {
		return []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s))
	}
	return []byte(s)
}

// Println 打印带有颜色格式的字符串，并换行
func (c Color) Println(s string) {
	fmt.Println(c.String(s))
}

// Printf 全文本颜色打印
// 格式化输入的字符串，并为整个结果添加颜色
func (c Color) Printf(s string, v ...any) {
	fmt.Print(c.String(fmt.Sprintf(s, v...)))
}

// Xprintf 仅变量文本颜色打印
// 仅为输入的变量部分添加颜色，格式字符串保持原样
func (c Color) Xprintf(s string, v ...any) {
	if len(v) > 0 {
		for i, a := range v {
			v[i] = c.String(fmt.Sprint(a))
		}
		fmt.Printf(s, v...)
	} else {
		fmt.Print(s)
	}
}

// Sprintf 全文本颜色格式化
// 格式化输入的字符串，并为整个结果添加颜色，返回格式化后的字符串
func (c Color) Sprintf(s string, v ...any) string {
	return c.String(fmt.Sprintf(s, v...))
}

// Xsprintf 仅变量文本颜色格式化
// 仅为输入的变量部分添加颜色，格式字符串保持原样，返回格式化后的字符串
func (c Color) Xsprintf(s string, v ...any) string {
	if len(v) > 0 {
		for i, a := range v {
			v[i] = c.String(fmt.Sprint(a))
		}
		return fmt.Sprintf(s, v...)
	} else {
		return fmt.Sprint(s)
	}
}
