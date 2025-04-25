package errorx

import (
	"fmt"
	"io"
	"runtime"
)

// Error 通用error
type Error struct {
	source error  // 源error
	msg    string // 报错信息
	stack  stack  // 调用栈
}

// 报错信息（用以实现error接口）
func (err *Error) Error() string {
	if err.source == nil {
		return err.msg
	} else {
		return err.msg + " | " + err.source.Error()
	}
}

// Unwrap 解包装
func (err *Error) Unwrap() error { return err.source }

// Format fmt打印实现
func (err *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		_, _ = fmt.Fprintf(s, "%v", err.Error())
		err.stack.Format(s, verb)
	case 's':
		_, _ = io.WriteString(s, err.Error())
	}
}

func New(v any) error {
	var err = &Error{stack: getStack()}
	switch e := v.(type) {
	case error:
		err.source = e
		err.msg = e.Error()
	case string:
		err.msg = e
	default:
		err.msg = fmt.Sprint(e)
	}
	return err
}

func Wrap(v any, msg string) error {
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

type stack []uintptr

// Format 打印调用栈信息（fmt实现）
func (s *stack) Format(f fmt.State, verb rune) {
	if verb == 'v' {
		i, frames := 1, runtime.CallersFrames(*s)
		for {
			if pc, more := frames.Next(); more && i <= 5 {
				_, _ = fmt.Fprintf(f, "\n%d : %s >> %s:%d", i, pc.Function, pc.File, pc.Line)
				i++
			} else {
				break
			}
		}
	}
}

// getStack 获取调用栈
func getStack() []uintptr {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[:n]
}
