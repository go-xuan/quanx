package errorx

import (
	"fmt"
	"io"
	"runtime"
)

// 通用error
type Error struct {
	source error
	msg    string
	stack  Stack
}

// 报错信息
func (err *Error) Error() string {
	if err.source == nil {
		return err.msg
	} else {
		return err.msg + ": " + err.source.Error()
	}
}

// 解包装
func (err *Error) Unwrap() error { return err.source }

// fmt打印实现
func (err *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		_, _ = fmt.Fprintf(s, "%v", err.Error())
		err.stack.Format(s, verb)
	case 's':
		_, _ = io.WriteString(s, err.Error())
	}
}

func New(v any) *Error {
	var err = &Error{stack: getStack()}
	switch e := v.(type) {
	case error:
		err.source = e
		err.msg = e.Error()
	case string:
		err.msg = e
	default:
		err.msg = fmt.Sprintf("%v", e)
	}
	return err
}

func Wrap(v any, msg string) *Error {
	var err = &Error{msg: msg}
	switch e := v.(type) {
	case *Error:
		err.source = e
		err.stack = e.stack
	case error:
		err.source = e
		err.stack = getStack()
	default:
		err.source = New(e)
		err.stack = getStack()
	}
	return err
}

func Unwrap(err error) error {
	if t, ok := err.(interface {
		Unwrap() error
	}); ok {
		return t.Unwrap()
	} else {
		return nil
	}
}

func Errorf(format string, a ...interface{}) error {
	return &Error{msg: fmt.Sprintf(format, a...), stack: getStack()}
}

func GetMessage(v any) string {
	switch i := v.(type) {
	case error:
		return i.Error()
	default:
		return fmt.Sprintf("%v", i)
	}
}

type Stack []uintptr

// fmt打印实现
func (s *Stack) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		frames := runtime.CallersFrames(*s)
		i := 1
		for {
			if pc, more := frames.Next(); more {
				_, _ = fmt.Fprintf(f, "\n%d : %s >> %s:%d", i, pc.Func.Name(), pc.File, pc.Line)
				i++
			} else {
				break
			}
		}
	}
}

func getStack() []uintptr {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[:n-1]
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
	return
}
