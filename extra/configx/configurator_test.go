package configx

import (
	"fmt"
	"testing"

	"github.com/go-xuan/quanx/types/anyx"
)

type test struct {
	Id   string `json:"id" default:"123"`
	Name string `json:"name" default:"test"`
}

func (t *test) Format() string {
	return "test show id"
}

func (t *test) Reader(from From) Reader {
	return nil
}

func (t *test) Execute() error {
	fmt.Println("==============test run logic start============")
	// 设置默认值
	if err := anyx.SetDefaultValue(t); err != nil {
		return err
	}
	// 自定义设置
	t.Name = "hello world"
	fmt.Println("==============test run logic end============")
	return nil
}

func TestConfigurator(t *testing.T) {
	var config = &test{}
	fmt.Println("before execute :", config)
	if err := config.Execute(); err != nil {
		t.Error(err)
		return
	}
	fmt.Println("after execute :", config)
}
