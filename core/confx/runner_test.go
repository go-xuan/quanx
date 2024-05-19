package confx

import (
	"fmt"
	"github.com/go-xuan/quanx/core"
	"testing"
)

type Test struct{}

func (t Test) Title() string {
	return "test show theme"
}

func (t Test) Reader() *Reader {
	return nil
}

func (t Test) Run() error {
	fmt.Println("test run")
	return nil
}

func TestConfigurator(t *testing.T) {
	var config = &Test{}
	var engine = core.GetEngine()
	engine.AddConfigurator(config)
	fmt.Println(config.Title())
	if err := config.Run(); err != nil {
		return
	}

}
