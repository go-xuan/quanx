package sqlx

import (
	"fmt"
	"github.com/quanxiaoxuan/quanx/common/constx"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/quanxiaoxuan/quanx/utils/timex"
)

const (
	paramKey     = "param_key"
	paramType    = "param_type"
	goTypeString = "string"
	goTypeInt    = "int"
	goTypeInt64  = "int64"
	tagTypeDate  = "date"
	tagTypeTime  = "time"
)

// 通过反射标签，构建完整sql(通用方法)
func SqlPrepared(sql string, obj interface{}) string {
	typeRef := reflect.TypeOf(obj)
	valueRef := reflect.ValueOf(obj)
	for i := 0; i < typeRef.NumField(); i++ {
		pKey := typeRef.Field(i).Tag.Get(paramKey)
		pType := typeRef.Field(i).Tag.Get(paramType)
		goType := typeRef.Field(i).Type.Name()
		pValue := valueRef.Field(i).String()
		switch goType {
		case goTypeString:
			pValue = `'` + valueRef.Field(i).String() + `'`
		case goTypeInt:
			pValue = strconv.FormatInt(valueRef.Field(i).Int(), 10)
		case goTypeInt64:
			switch pType {
			case tagTypeTime:
				pValue = `'` + timex.SecondFormat(valueRef.Field(i).Int()/1000, constx.TimeFmt) + `'`
			case tagTypeDate:
				pValue = `'` + timex.SecondFormat(valueRef.Field(i).Int()/1000, constx.DateFmt) + `'`
			default:
				pValue = strconv.FormatInt(valueRef.Field(i).Int(), 10)
			}
		}
		sql = strings.ReplaceAll(sql, pKey, pValue)
	}
	return sql
}

type Model struct {
	Id   int    `json:"id" param_key:"#{id}"`
	Name string `json:"name" param_key:"#{name}"`
	Num  int64  `json:"num" param_key:"#{num}"`
	Date int64  `json:"date" param_key:"#{date}" param_type:"date"`
	Time int64  `json:"time" param_key:"#{time}" param_type:"time"`
}

// 测试demo
func demo() {
	sql := "select #{id},#{num},#{name},#{date}，#{time}"
	var model = Model{
		Id:   1,
		Name: "Abc",
		Num:  12356879,
		Date: 1669951433000,
		Time: time.Now().UnixMilli(),
	}
	sql = SqlPrepared(sql, model)
	fmt.Println(sql)
}
