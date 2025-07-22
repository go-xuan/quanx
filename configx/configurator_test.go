package configx

import (
	"fmt"
	"testing"
)

type test struct {
	Id   string `json:"id" default:"123"`
	Name string `json:"name" default:"test"`
}

func (t *test) NeedRead() bool {
	if t.Id == "" && t.Name == "" {
		return true
	}
	return false
}

func (t *test) Reader(from From) Reader {
	return nil
}

func (t *test) Execute() error {
	t.Name = "hello world"
	return nil
}

func TestConfigurator(t *testing.T) {
	var config = &test{}
	fmt.Println("before execute :", config)
	if err := ReadAndExecute(config, FromTag); err != nil {
		t.Error(err)
		return
	}
	fmt.Println("after execute :", config)
}
