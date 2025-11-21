package configx

import (
	"fmt"
	"testing"
)

type test struct {
	Id   string `json:"id" default:"123"`
	Name string `json:"name" default:"test"`
	Desc string `json:"desc" default:""`
}

func (t *test) Valid() bool {
	return t.Id != "" && t.Name != ""
}

func (t *test) Readers() []Reader {
	return []Reader{
		NewTagReader(),
	}
}

func (t *test) Execute() error {
	t.Desc = "hello world"
	return nil
}

func TestConfigurator(t *testing.T) {
	var config = &test{}
	fmt.Println("before execute :", config)
	if err := LoadConfigurator(config); err != nil {
		t.Error(err)
		return
	}
	fmt.Println("after execute :", config)
}
