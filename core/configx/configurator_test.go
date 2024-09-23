package configx

import (
	"fmt"
	"github.com/go-xuan/quanx"
	"testing"
)

type Test struct{}

func (t Test) Format() string {
	return ""
}

func (t Test) ID() string {
	return "test show theme"
}

func (t Test) Reader() *Reader {
	return nil
}

func (t Test) Execute() error {
	fmt.Println("test run")
	return nil
}

func TestConfigurator(t *testing.T) {
	var config = &Test{}
	var engine = quanx.NewEngine()
	engine.AddConfigurator(config)
	fmt.Println(config.ID())
	if err := config.Execute(); err != nil {
		return
	}
}
