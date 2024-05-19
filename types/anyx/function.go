package anyx

import "github.com/go-xuan/quanx/os/errorx"

func Return[T any](t T, err error) T {
	errorx.Panic(err)
	return t
}
