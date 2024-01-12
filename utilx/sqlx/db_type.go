package sqlx

// 共用数据类型
const (
	String     = "string"
	Varchar    = "varchar"
	Char       = "char"
	Varchar255 = "varchar(255)"
	Text       = "text"
	Tinyint    = "tinyint"
	Smallint   = "smallint"
	Mediumint  = "mediumint"
	Int        = "int"
	Bigint     = "bigint"
	Int2       = "int2"
	Int4       = "int4"
	Int8       = "int8"
	Int64      = "int64"
	Float4     = "float4"
	Numeric    = "numeric" // 数字
	Numeric6   = "numeric(10,6)"
	Numeric2   = "numeric(10,2)"
	Decimal    = "decimal"
	Time       = "time"
	GoTime     = "time.Time"
	Timestamp  = "timestamp"
	Date       = "date"
	Datetime   = "datetime"
	Bool       = "bool"
)

// java类常用类型
const (
	JavaString     = "String"
	JavaInteger    = "Integer"
	JavaInt        = "int"
	JavaLong       = "Long"
	JavaDate       = "Date"
	JavaBigDecimal = "BigDecimal"
	JavaFloat      = "Float"
	JavaBoolean    = "Boolean"
)
