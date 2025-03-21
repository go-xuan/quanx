package configx

import (
	"fmt"
	"testing"

	"github.com/go-xuan/quanx/types/anyx"
)

type Test struct {
	Id   string `json:"id" default:"123"`
	Name string `json:"name" default:"test"`
}

func (t *Test) Format() string {
	return "test show id"
}

func (t *Test) Reader(from From) Reader {
	return nil
}

func (t *Test) Execute() error {
	// 设置默认值
	if err := anyx.SetDefaultValue(t); err != nil {
		return err
	}
	fmt.Println("test run logic start")
	// 自定义设置
	t.Name = "hello world"
	fmt.Println("test run logic end")
	return nil
}

func TestConfigurator(t *testing.T) {
	var config = &Test{}
	fmt.Println("before execute, id:", config.Id)
	if err := config.Execute(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("==========================")
	fmt.Println("after execute, id:", config.Id)
	fmt.Println(config)
}
