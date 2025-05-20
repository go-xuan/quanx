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
	v := floatx.NewFloat64(123.4)
	fmt.Println("string = ", v.String())
	fmt.Println("int = ", v.Int(111))
	fmt.Println("int64 = ", v.Int64())
	fmt.Println("float64 = ", v.Float64())
	fmt.Println("bool = ", v.Bool())
	fmt.Println("string = ", v.String())
}

func TestMarshal(t *testing.T) {
	type Demo struct {
		String stringx.String `json:"name"`
		Time   timex.Time     `json:"create_time"`
		Date   timex.Date     `json:"create_date"`
		Bool   boolx.Bool     `json:"bool"`
		Int    intx.Int       `json:"int"`
		Int64  intx.Int64     `json:"int64"`
		Float  floatx.Float   `json:"float"`
	}

	bytes := []byte(`{"name":null,"create_time": 11111,"create_date":"2024-11-21","bool":1,"int":47826,"int64":23364,"float":57575.138063,"value":123.4}`)
	demo := &Demo{}
	if err := json.Unmarshal(bytes, demo); err != nil {
		panic(err)
	}
	bytes, _ = json.Marshal(demo)
	fmt.Println("反序列化：", string(bytes))

	demo.String = stringx.NewString(randx.String())
	demo.Bool = boolx.NewBool(randx.Bool())
	demo.Date = timex.NewDate(randx.Time())
	demo.Time = timex.NewTime(randx.Time())
	demo.Int = intx.NewInt(randx.Int())
	demo.Int64 = intx.NewInt64(randx.Int64())
	demo.Float = floatx.NewFloat64(randx.Float64())
	bytes, _ = json.Marshal(demo)
	fmt.Println("重新赋值：", string(bytes))
}
