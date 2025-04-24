package maskx

import (
	"fmt"
	"testing"

	"github.com/go-xuan/quanx/utils/randx"
)

func TestDesensitize(t *testing.T) {
	phone := randx.Phone()
	name := randx.Name()
	idCard := randx.IdCard()
	email := randx.Email()
	fmt.Println(phone, "==>", Phone.Desensitize(phone))
	fmt.Println(name, "==>", Name.Desensitize(name))
	fmt.Println(idCard, "==>", IdCard.Desensitize(idCard))
	fmt.Println(email, "==>", Email.Desensitize(email))
}
