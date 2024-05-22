package fmtx

import "fmt"

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
