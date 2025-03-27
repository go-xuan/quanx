package enumx

import (
	"fmt"
	"testing"
)

func TestEnum(t *testing.T) {
	// 声明key为string类型的枚举，两种声明方式等价
	var sf = NewStringEnum[float64]()    // *Enum[string, float64]
	var sf2 = NewEnum[string, float64]() // *Enum[string, float64]

	// 声明key为int类型的枚举类，两种声明方式等价
	var is = NewIntEnum[string]()    // *Enum[int, string]
	var is2 = NewEnum[int, string]() // *Enum[int, string]

	// 声明key为任意comparable类型的枚举类，按需声明
	sa := NewEnum[string, any]()

	sf.Add("1", 111.111).Add("2", 222.222).Add("3", 333.333)
	sf2.Add("1", 111.111).Add("2", 222.222).Add("3", 333.333)

	is.Add(1, "AAA").Add(2, "BBB").Add(3, "CCC")
	is2.Add(1, "AAA").Add(2, "BBB").Add(3, "CCC")

	sa.Add("1", 1111).Add("2", 2222).Add("3", 3333)

	fmt.Println(sf.Get("1") == sf2.Get("1"))
	fmt.Println(is.Get(1) == is2.Get(1))
	fmt.Println(sa.Keys())
	sa.Remove("1")
	fmt.Println(sa.Keys())

	fmt.Println(sf.Values())
	sf.Clear()
	fmt.Println(sf.Values())
}
