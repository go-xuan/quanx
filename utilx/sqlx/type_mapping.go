package sqlx

var (
	db2Go   map[string]string
	db2Gorm map[string]string
	db2Java map[string]string
)

// DB-Go类型映射
func DB2Go() map[string]string {
	if db2Go == nil {
		db2Go = make(map[string]string)
		db2Go[Char] = String
		db2Go[Varchar] = String
		db2Go[Text] = String
		db2Go[Int2] = Int
		db2Go[Int4] = Int
		db2Go[Int8] = Int64
		db2Go[Tinyint] = Int
		db2Go[Smallint] = Int
		db2Go[Mediumint] = Int
		db2Go[Int] = Int
		db2Go[Bigint] = Int64
		db2Go[Float4] = Float4
		db2Go[Numeric] = Float4
		db2Go[Timestamp] = GoTime
		db2Go[Datetime] = GoTime
		db2Go[Date] = GoTime
		db2Go[Bool] = Bool
	}
	return db2Go
}

// DB-Gorm类型映射
func DB2Gorm() map[string]string {
	if db2Gorm == nil {
		db2Gorm = make(map[string]string)
		db2Gorm[Char] = Char
		db2Gorm[Varchar] = Varchar255
		db2Gorm[Text] = Text
		db2Gorm[Int2] = Smallint
		db2Gorm[Int4] = Int
		db2Gorm[Int8] = Int64
		db2Gorm[Tinyint] = Tinyint
		db2Gorm[Smallint] = Smallint
		db2Gorm[Mediumint] = Mediumint
		db2Gorm[Int] = Int
		db2Gorm[Bigint] = Bigint
		db2Gorm[Float4] = Numeric6
		db2Gorm[Numeric] = Numeric2
		db2Gorm[Timestamp] = Timestamp
		db2Gorm[Datetime] = Timestamp
		db2Gorm[Date] = Date
		db2Gorm[Bool] = Bool
	}
	return db2Gorm
}

// DB-java类型映射
func DB2Java() map[string]string {
	if db2Java == nil {
		db2Java = make(map[string]string)
		db2Java[Char] = JavaString
		db2Java[Varchar] = JavaString
		db2Java[Text] = JavaString
		db2Java[Int2] = JavaInteger
		db2Java[Int4] = JavaLong
		db2Java[Int8] = JavaLong
		db2Java[Tinyint] = JavaInt
		db2Java[Smallint] = JavaInteger
		db2Java[Mediumint] = JavaInteger
		db2Java[Int] = JavaLong
		db2Java[Bigint] = JavaLong
		db2Java[Float4] = JavaFloat
		db2Java[Numeric] = JavaFloat
		db2Java[Decimal] = JavaBigDecimal
		db2Java[Timestamp] = JavaDate
		db2Java[Datetime] = JavaDate
		db2Java[Date] = JavaDate
		db2Java[Bool] = JavaBoolean
	}
	return db2Java
}
