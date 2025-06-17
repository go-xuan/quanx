package wrapx

import (
	"context"

	"github.com/go-xuan/quanx/base/errorx"
)

func NewWrapper() *Wrapper {
	return &Wrapper{
		wraps: make([]FuncWrap, 0),
	}
}

// Func 函数
type Func func(ctx context.Context) error

// FuncWrap 函数包装器
type FuncWrap func(fn Func) Func

// Wrapper 包装器
type Wrapper struct {
	wraps []FuncWrap
}

// Add 添加包装器
func (w *Wrapper) Add(wraps ...FuncWrap) {
	if len(wraps) > 0 {
		w.wraps = append(w.wraps, wraps...)
	}
}

// Call 调用
func (w *Wrapper) Call(ctx context.Context, fn Func) error {
	for _, wrapper := range w.wraps {
		fn = wrapper(fn)
	}
	if err := fn(ctx); err != nil {
		return errorx.Wrap(err, "wrapper call error")
	}
	return nil
}
