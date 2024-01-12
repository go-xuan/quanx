package sqlx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-xuan/quanx/utilx/timex"
)

// 通过反射标签，构建完整sql(通用方法)
func SqlPrepared(sql string, obj interface{}) string {
	typeRef := reflect.TypeOf(obj)
	valueRef := reflect.ValueOf(obj)
	for i := 0; i < typeRef.NumField(); i++ {
		pKey := typeRef.Field(i).Tag.Get("param_key")
		pType := typeRef.Field(i).Tag.Get("param_type")
		goType := typeRef.Field(i).Type.Name()
		pValue := valueRef.Field(i).String()
		switch goType {
		case String:
			pValue = `'` + valueRef.Field(i).String() + `'`
		case Int:
			pValue = strconv.FormatInt(valueRef.Field(i).Int(), 10)
		case Int64:
			switch pType {
			case Time:
				pValue = `'` + timex.SecondFormat(valueRef.Field(i).Int()/1000, timex.TimeFmt) + `'`
			case Date:
				pValue = `'` + timex.SecondFormat(valueRef.Field(i).Int()/1000, timex.DateFmt) + `'`
			default:
				pValue = strconv.FormatInt(valueRef.Field(i).Int(), 10)
			}
		}
		sql = strings.ReplaceAll(sql, pKey, pValue)
	}
	return sql
}

type Model struct {
	Id   int       `json:"id" param_key:"#{id}"`
	Name string    `json:"name" param_key:"#{name}"`
	Num  int64     `json:"num" param_key:"#{num}"`
	Date time.Time `json:"date" param_key:"#{date}" param_type:"date"`
	Time time.Time `json:"time" param_key:"#{time}" param_type:"time"`
}

// 测试demo
func demo() {
	sql := "select #{id},#{num},#{name},#{date}，#{time}"
	var model = Model{
		Id:   1,
		Name: "Abc",
		Num:  12356879,
		Date: time.Now(),
		Time: time.Now(),
	}
	sql = SqlPrepared(sql, model)
	fmt.Println(sql)
}
