package hugegraphx

// 常量配置
const (
	Propertykeys = "propertykeys" // 属性API
	Vertexlabels = "vertexlabels" // 顶点API
	Edgelabels   = "edgelabels"   // 边API
	Indexlabels  = "indexlabels"  // 边API
)

const (
	Multiple       = "MULTIPLE"                                                                       // 允许多次连接
	Single         = "SINGLE"                                                                         // 单条连接
	DateTypesAll   = "TEXT,BYTE,DATE,INT,LONG,BOOLEAN,DOUBLE,FLOAT,UUID,BLOB,OBJECT,UNKNOWN"          // 属性数据类型
	CardinalityAll = "SINGLE,SET,LIST"                                                                // 属性类型基数
	IdStrategyAll  = "DEFAULT,AUTOMATIC,PRIMARY_KEY,CUSTOMIZE_STRING,CUSTOMIZE_NUMBER,CUSTOMIZE_UUID" // ID策略
	ThicknessAll   = "THICK,NORMAL,FINE"                                                              // 边线条粗细
	SizeAll        = "TINY,SMALL,NORMAL,BIG,HUGE"                                                     // 顶点样式大小
	BaseTypeAll    = "VERTEX_LABEL,EDGE_LABEL"                                                        // 模型类型
	IndexTypeAll   = "SECONDARY,RANGE,SEARCH"                                                         // 索引类型
)
