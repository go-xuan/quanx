package anyx

import (
	"fmt"
	"testing"
)

func TestValue(t *testing.T) {
	v := ValueOf("123.67")
	fmt.Printf("%t\n", v)
	fmt.Println("int = ", v.Int())
	fmt.Println("int64 = ", v.Int64())
	fmt.Println("float64 = ", v.Float64())
	fmt.Println("bool = ", v.Bool())
	fmt.Println("string = ", v.String())
}
