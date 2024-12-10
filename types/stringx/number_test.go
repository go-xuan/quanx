package stringx

import (
	"fmt"
	"testing"
)

func TestNumber(t *testing.T) {
	fmt.Println(ConvertArabToChinese(123))
	fmt.Println(ConvertArabToChinese(1234))
	fmt.Println(ConvertArabToChinese(12345))
	fmt.Println(ConvertArabToChinese(123456))
	fmt.Println(ConvertArabToChinese(1234567))
	fmt.Println(ConvertArabToChinese(12345678))
	fmt.Println(ConvertArabToChinese(123456789))
}
