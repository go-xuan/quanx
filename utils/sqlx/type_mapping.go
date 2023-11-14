package sqlx

var Pg2GoTypeMap = initPg2GoTypeMap()
var Pg2GormTypeMap = initPg2GormTypeMap()
var Pg2JavaTypeMap = initPg2JavaTypeMap()
var Pg2CkTypeMap = initPg2CkTypeMap()

// PG-Go类型映射
func initPg2GoTypeMap() map[string]string {
	var typeMap = make(map[string]string)
	typeMap[Varchar] = String
	typeMap[Text] = String
	typeMap[Int2] = Int
	typeMap[Int4] = Int
	typeMap[Int8] = Int64
	typeMap[Timestamp] = Time
	typeMap[Date] = Time
	typeMap[Float4] = Float4
	typeMap[Numeric] = Float4
	typeMap[Bool] = Bool
	return typeMap
}

// mysql-Go类型映射
func initMysql2GoTypeMap() map[string]string {
	var typeMap = make(map[string]string)
	typeMap[Varchar] = String
	typeMap[Text] = String
	typeMap[Tinyint] = Int
	typeMap[Int4] = Int
	typeMap[Int8] = Int64
	typeMap[Timestamp] = Time
	typeMap[Date] = Time
	typeMap[Float4] = Float4
	typeMap[Numeric] = Float4
	typeMap[Bool] = Bool
	return typeMap
}

// PG-Gorm类型映射
func initPg2GormTypeMap() map[string]string {
	var typeMap = make(map[string]string)
	typeMap[Varchar] = Varchar255
	typeMap[Text] = Text
	typeMap[Int2] = Smallint
	typeMap[Int4] = Int
	typeMap[Int8] = Bigint
	typeMap[Timestamp] = Timestamp
	typeMap[Date] = Date
	typeMap[Float4] = Numeric6
	typeMap[Numeric] = Numeric2
	typeMap[Bool] = Bool
	return typeMap
}

// PG-java类型映射
func initPg2JavaTypeMap() map[string]string {
	var typeMap = make(map[string]string)
	typeMap[Varchar] = JavaString
	typeMap[Text] = JavaString
	typeMap[Int2] = JavaInteger
	typeMap[Int4] = JavaLong
	typeMap[Int8] = JavaLong
	typeMap[Timestamp] = JavaDate
	typeMap[Date] = JavaDate
	typeMap[Float4] = JavaBigDecimal
	typeMap[Numeric] = JavaBigDecimal
	typeMap[Bool] = JavaBoolean
	return typeMap
}

// PG-ClickHouse类型映射
func initPg2CkTypeMap() map[string]string {
	var typeMap = make(map[string]string)
	typeMap[Varchar] = CkString
	typeMap[Text] = CkString
	typeMap[Int2] = CkInt8
	typeMap[Int4] = CkInt16
	typeMap[Int8] = CkInt32
	typeMap[Timestamp] = CkDateTime
	typeMap[Date] = CkDate
	typeMap[Float4] = CkFloat64
	typeMap[Numeric] = CkFloat64
	typeMap[Bool] = CkBool
	return typeMap
}
