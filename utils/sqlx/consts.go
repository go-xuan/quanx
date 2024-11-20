package sqlx

// symbol
const (
	Empty         = ""
	Blank         = " "
	LeftBracket   = "("
	RightBracket  = ")"
	Comma         = ","
	NewLine       = "\n"
	Equals        = "="
	ReplacePrefix = "value@"
)

// keyword
const (
	SELECT    = "select"
	UPDATE    = "update"
	CREATE    = "create"
	DELETE    = "delete"
	INSERT    = "insert"
	FROM      = "from"
	WHERE     = "where"
	SET       = "set"
	LEFT      = "left"
	RIGHT     = "right"
	INNER     = "inner"
	OUTER     = "outer"
	JOIN      = "join"
	GROUP     = "group"
	GroupBy   = "group by"
	ORDER     = "order"
	OrderBy   = "order by"
	HAVING    = "having"
	LIMIT     = "limit"
	OFFSET    = "offset"
	AS        = "as"
	AND       = "and"
	ON        = "on"
	OR        = "or"
	IN        = "in"
	NOT       = "not"
	LIKE      = "like"
	By        = "by"
	DISTINCT  = "distinct"
	OVER      = "over"
	PARTITION = "partition"
	CASE      = "case"
	WHEN      = "when"
	THEN      = "then"
	END       = "end"
	ASC       = "asc"
	DESC      = "desc"
)

// 随机数生成器-数据类型
const (
	String     = "string"       // 字符串
	Text       = "text"         // 文本
	Varchar    = "varchar"      // 字符
	Varchar100 = "varchar(100)" // 100字符串
	Char       = "char"         // 字节
	Int        = "int"          // 数字
	Int2       = "int2"         // 小整型
	Int4       = "int4"         // 中整型
	Int8       = "int8"         // 大整型
	Int64      = "int64"        // 64位整数
	Float      = "float"        // 浮点数
	Float4     = "float4"       // 浮点
	Float64    = "float64"      // 浮点数
	Tinyint    = "tinyint"      // 微整数
	Smallint   = "smallint"     // 小整数
	Mediumint  = "mediumint"    // 中整数
	Bigint     = "bigint"       // 大整数
	Sequence   = "sequence"     // 序列
	Uuid       = "uuid"         // UUID
	Bool       = "bool"         // 布尔
	Numeric    = "numeric"      // 数字
	Time       = "time"         // 时间
	Date       = "date"         // 日期
	Timestamp  = "timestamp"    // 时间戳
	Timestampz = "timestamptz"  // 时间戳带时区
	Datetime   = "datetime"     // 日期时间
	TimeTime   = "time.Time"    // 时间
)
