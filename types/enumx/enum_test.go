package enumx

import (
	"fmt"
	"testing"
)

func TestEnum(t *testing.T) {
	ss := NewStringEnum[string]().
		Add("1", "1111").
		Add("2", "2222").
		Add("3", "3333")
	fmt.Println(ss.Keys())
	fmt.Println(ss.Values())
	fmt.Println(ss.Get("1"))

	ii := NewIntEnum[int]().
		Add(1, 1111).
		Add(2, 2222).
		Add(3, 3333)
	fmt.Println(ii.Keys())
	fmt.Println(ii.Values())
	fmt.Println(ii.Get(1))

	is := NewIntEnum[string]().
		Add(1, "1111").
		Add(2, "2222").
		Add(3, "3333")
	fmt.Println(is.Keys())
	fmt.Println(is.Values())
	fmt.Println(is.Get(1))

	si := NewStringEnum[int]().
		Add("1", 1111).
		Add("2", 2222).
		Add("3", 3333)
	fmt.Println(si.Keys())
	fmt.Println(si.Values())
	fmt.Println(si.Get("1"))
}
