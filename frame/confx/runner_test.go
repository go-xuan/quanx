package confx

import (
	"fmt"
	"github.com/go-xuan/quanx"
	"testing"
)

type Test struct{}

func (t Test) Theme() string {
	return "test show theme"
}

func (t Test) Reader() *Reader {
	return nil
}

func (t Test) Run() error {
	return nil
}

func TestConfigurator(t *testing.T) {
	var config = &Test{}
	var engine = quanx.GetEngine()
	engine.AddConfigurator(config)
	fmt.Println(config.Theme())
	if err := config.Run(); err != nil {
		return
	}

}
