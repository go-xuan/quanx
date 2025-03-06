package anyx

import (
	"fmt"
	"testing"
)

func TestStruct(t *testing.T) {
	var demo = &struct {
		Name  string  `json:"name" default:"abc"`
		Num   string  `json:"num" default:"123"`
		Float float64 `json:"float" default:"123.45"`
		Null  string  `json:"null"`
	}{}
	if err := SetDefaultValue(demo); err != nil {
		t.Fatal(err)
	}
	fmt.Println("name = ", demo.Name)
	fmt.Println("num = ", demo.Num)
	fmt.Println("float = ", demo.Float)
	fmt.Println("null = ", demo.Null)
}
