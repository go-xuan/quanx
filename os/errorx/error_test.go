package errorx

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err := errors.New("test error")
	fmt.Println("原生error：", err)

	err = Errorf("errorx：%v", err)
	fmt.Println(err)

	// err封包
	err = Wrap(err, "第一层包装")
	fmt.Println(err)
	err = Wrap(err, "第二层包装")
	fmt.Println(err)
	err = Wrap(err, "第三层包装")
	fmt.Println(err)

	// error解包
	err = Unwrap(err)
	fmt.Println(err)
	// error解包
	err = Unwrap(err)
	fmt.Println(err)
}
