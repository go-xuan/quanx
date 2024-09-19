package configx

import (
	"fmt"
	"testing"

	"github.com/go-xuan/quanx/app"
)

type Test struct{}

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
	var engine = app.NewEngine()
	engine.AddConfigurator(config)
	fmt.Println(config.ID())
	if err := config.Execute(); err != nil {
		return
	}

}
