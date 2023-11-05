package sqlx

import "github.com/go-xuan/quanx/common/constx"

var Pg2GoTypeMap = initPg2GoTypeMap()
var Pg2GormTypeMap = initPg2GormTypeMap()
var Pg2JavaTypeMap = initPg2JavaTypeMap()
var Pg2CkTypeMap = initPg2CkTypeMap()

// PG-Go类型映射
func initPg2GoTypeMap() map[string]string {
	var typeMap = make(map[string]string)
	typeMap[constx.Varchar] = constx.String
	typeMap[constx.Text] = constx.String
	typeMap[constx.Int2] = constx.Int
	typeMap[constx.Int4] = constx.Int
	typeMap[constx.Int8] = constx.Int64
	typeMap[constx.Timestamp] = constx.Time
	typeMap[constx.Date] = constx.Time
	typeMap[constx.Float4] = constx.Float4
	typeMap[constx.Numeric] = constx.Float4
	typeMap[constx.Bool] = constx.Bool
	return typeMap
}

// mysql-Go类型映射
func initMysql2GoTypeMap() map[string]string {
	var typeMap = make(map[string]string)
	typeMap[constx.Varchar] = constx.String
	typeMap[constx.Text] = constx.String
	typeMap[constx.Tinyint] = constx.Int
	typeMap[constx.Int4] = constx.Int
	typeMap[constx.Int8] = constx.Int64
	typeMap[constx.Timestamp] = constx.Time
	typeMap[constx.Date] = constx.Time
	typeMap[constx.Float4] = constx.Float4
	typeMap[constx.Numeric] = constx.Float4
	typeMap[constx.Bool] = constx.Bool
	return typeMap
}

// PG-Gorm类型映射
func initPg2GormTypeMap() map[string]string {
	var typeMap = make(map[string]string)
	typeMap[constx.Varchar] = constx.Varchar255
	typeMap[constx.Text] = constx.Text
	typeMap[constx.Int2] = constx.Smallint
	typeMap[constx.Int4] = constx.Int
	typeMap[constx.Int8] = constx.Bigint
	typeMap[constx.Timestamp] = constx.Timestamp
	typeMap[constx.Date] = constx.Date
	typeMap[constx.Float4] = constx.Numeric6
	typeMap[constx.Numeric] = constx.Numeric2
	typeMap[constx.Bool] = constx.Bool
	return typeMap
}

// PG-java类型映射
func initPg2JavaTypeMap() map[string]string {
	var typeMap = make(map[string]string)
	typeMap[constx.Varchar] = constx.JavaString
	typeMap[constx.Text] = constx.JavaString
	typeMap[constx.Int2] = constx.JavaInteger
	typeMap[constx.Int4] = constx.JavaLong
	typeMap[constx.Int8] = constx.JavaLong
	typeMap[constx.Timestamp] = constx.JavaDate
	typeMap[constx.Date] = constx.JavaDate
	typeMap[constx.Float4] = constx.JavaBigDecimal
	typeMap[constx.Numeric] = constx.JavaBigDecimal
	typeMap[constx.Bool] = constx.JavaBoolean
	return typeMap
}

// PG-ClickHouse类型映射
func initPg2CkTypeMap() map[string]string {
	var typeMap = make(map[string]string)
	typeMap[constx.Varchar] = constx.CkString
	typeMap[constx.Text] = constx.CkString
	typeMap[constx.Int2] = constx.CkInt8
	typeMap[constx.Int4] = constx.CkInt16
	typeMap[constx.Int8] = constx.CkInt32
	typeMap[constx.Timestamp] = constx.CkDateTime
	typeMap[constx.Date] = constx.CkDate
	typeMap[constx.Float4] = constx.CkFloat64
	typeMap[constx.Numeric] = constx.CkFloat64
	typeMap[constx.Bool] = constx.CkBool
	return typeMap
}
