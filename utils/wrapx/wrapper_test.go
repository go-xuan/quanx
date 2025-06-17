package wrapx

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-xuan/quanx/base/contextx"
)

func TestWrapper(t *testing.T) {
	wrapper := NewWrapper()
	wrapper.Add(setValue, cost)

	i := 1
	fn := func(ctx context.Context) error {
		i += 10
		fmt.Println("get value:", contextx.GetValue(ctx, "quanx").Int())
		return nil
	}

	ctx := contextx.New()
	contextx.SetValue(ctx, "quanx", i)

	fmt.Println("before call: ", i)
	if err := wrapper.Call(ctx, fn); err != nil {
		t.Error(err)
	}
	fmt.Println("after  call: ", i)
}

func cost(fn Func) Func {
	return func(ctx context.Context) error {
		fmt.Println("start cost")
		now := time.Now()
		if err := fn(ctx); err != nil {
			return err
		}
		fmt.Println("cost: ", time.Since(now))
		return nil
	}
}

func setValue(fn Func) Func {
	return func(ctx context.Context) error {
		before := contextx.GetValue(ctx, "quanx").Int()
		after := before + 3
		contextx.SetValue(ctx, "quanx", after)
		fmt.Printf("set value: %v ==> %v \n", before, after)
		if err := fn(ctx); err != nil {
			return err
		}
		return nil
	}
}
