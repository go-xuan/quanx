package sqlx

// symbol
const (
	Empty         = ""
	Blank         = " "
	LeftBracket   = "("
	RightBracket  = ")"
	Comma         = ","
	NewLine       = "\n"
	ReplacePrefix = "value@"
)

// keyword
const (
	SELECT    = "select"
	FROM      = "from"
	WHERE     = "where"
	LEFT      = "left"
	RIGHT     = "right"
	INNER     = "inner"
	JOIN      = "join"
	GROUP     = "group"
	ORDER     = "order"
	HAVING    = "having"
	LIMIT     = "limit"
	OFFSET    = "offset"
	AS        = "as"
	AND       = "and"
	ON        = "on"
	IN        = "in"
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

var KEYWORDS = []string{
	SELECT, FROM, WHERE, JOIN, GROUP, ORDER, HAVING, LIMIT, OFFSET,
	ASC, DESC, CASE, WHEN, THEN, END, INNER, LEFT, RIGHT,
	DISTINCT, PARTITION, OVER, AS, AND, ON, IN, By,
}
