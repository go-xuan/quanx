package anyx

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-xuan/quanx/types/boolx"
	"github.com/go-xuan/quanx/types/floatx"
	"github.com/go-xuan/quanx/types/intx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/types/timex"
	"github.com/go-xuan/quanx/utils/randx"
)

func TestValue(t *testing.T) {
	type Demo struct {
		String stringx.String `json:"name"`
		Time   timex.Time     `json:"create_time"`
		Date   timex.Date     `json:"create_date"`
		Bool   boolx.Bool     `json:"bool"`
		Int    intx.Int       `json:"int"`
		Int64  intx.Int64     `json:"int64"`
		Float  floatx.Float   `json:"float"`
		//Value  Value          `json:"value"`
	}

	j := `{"name":12345678,"time":2004565434,"date":"2024-11-21","bool":1,"int":47826,"int64":23364,"float":57575.138063,"value":123.4}`
	base := &Demo{}
	if err := json.Unmarshal([]byte(j), base); err != nil {
		panic(err)
	}
	b, _ := json.Marshal(base)
	fmt.Println(string(b))

	base.String = stringx.NewString(randx.String())
	base.Bool = boolx.NewBool(randx.Bool())
	base.Date = timex.NewDate(randx.Time())
	base.Time = timex.NewTime(randx.Time())
	base.Int = intx.NewInt(randx.Int())
	base.Int64 = intx.NewInt64(randx.Int64())
	base.Float = floatx.NewFloat64(randx.Float64())

	b, _ = json.Marshal(base)
	fmt.Println(string(b))

	v := base.Float
	fmt.Println("string = ", v.String())
	fmt.Println("int = ", v.Int(111))
	fmt.Println("int64 = ", v.Int64())
	fmt.Println("float64 = ", v.Float64())
	fmt.Println("bool = ", v.Bool())
	fmt.Println("string = ", v.String())
}
