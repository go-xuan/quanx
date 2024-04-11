package hugegraphx

// post接口请求参数
type Param struct {
	Gremlin  string      `json:"gremlin" `  // gremlin执行语句
	Bindings interface{} `json:"bindings" ` // 绑定参数
	Language string      `json:"language" ` // 语言
	Aliases  interface{} `json:"aliases" `  // 别名
}

// post接口返回结果
type ApiResp[T any] struct {
	RequestId string    `json:"requestId"` // 请求ID
	Status    Status    `json:"status"`    // 返回状态
	Result    Result[T] `json:"result"`    // 返回结果
}

// PING接口返回结果
type PingResp struct {
	Versions *Versions `json:"versions"` // 请求ID
}

type Versions struct {
	Version string `json:"version"`
	Core    string `json:"core"`
	Gremlin string `json:"gremlin"`
	Api     string `json:"api"`
}

// 返回结果
type Result[T any] struct {
	Data T           `json:"data"` // 结果体
	Meta interface{} `json:"meta"` // 元数据
}

// 返回状态
type Status struct {
	Message    string      `json:"message"`    // 请求ID
	Code       int64       `json:"code"`       // 状态码
	Attributes interface{} `json:"attributes"` // 属性
}

// hugegraph查询【顶点】返回的data结果
type Vertexs[T any] []*Vertex[T]
type Vertex[T any] struct {
	Id         string `json:"id"`         // 主键id
	Label      string `json:"label"`      // 类型
	Type       string `json:"type"`       // 分类
	Properties T      `json:"properties"` // 属性
}

// hugegraph查询【边】返回的data结果
type Edges[T any] []*Edge[T]
type Edge[T any] struct {
	Id         string `json:"id"`         // ID
	Label      string `json:"label"`      // 类型
	Type       string `json:"type"`       // 分类
	OutV       string `json:"outV"`       // 出顶点ID
	OutVLabel  string `json:"outVLabel"`  // 出顶点类型
	InV        string `json:"inV"`        // 类入顶点ID型
	InVLabel   string `json:"inVLabel"`   // 入顶点类型
	Properties T      `json:"properties"` // 属性
}

// hugegraph查询path()时返回的data结果
type Paths[T any] []*Path[T]
type Path[T any] struct {
	Labels  interface{}    `json:"labels"`
	Objects PathObjects[T] `json:"objects"`
}

// path()查询结果的【路径】中节点对象，包含顶点和边
type PathObjects[T any] []*PathObject[T]
type PathObject[T any] struct {
	Id         string `json:"id"`         // ID
	Label      string `json:"label"`      // 类型
	Type       string `json:"type"`       // 分类
	OutV       string `json:"outV"`       // 出顶点ID
	OutVLabel  string `json:"outVLabel"`  // 出顶点类型
	InV        string `json:"inV"`        // 类入顶点ID型
	InVLabel   string `json:"inVLabel"`   // 入顶点类型
	Properties T      `json:"properties"` // 对象属性
}

// 新增属性参数
type PropertyAdd struct {
	Name        string `json:"name"`        // 属性名称
	DataType    string `json:"data_type"`   // 属性类型
	Cardinality string `json:"cardinality"` // 属性类型基数
}

// 新增顶点参数
type VertexAdd struct {
	Name             string            `json:"name"`               // 顶点名称
	IdStrategy       string            `json:"id_strategy"`        // 主键策略
	Properties       []string          `json:"properties"`         // 属性列表
	PrimaryKeys      []string          `json:"primary_keys"`       // 主键属性列表
	NullableKeys     []string          `json:"nullable_keys"`      // 可空属性列表
	IndexLabels      []string          `json:"index_labels"`       // 索引列表
	Ttl              int               `json:"ttl"`                // TTL
	EnableLabelIndex bool              `json:"enable_label_index"` // 启用类型索引,默认为true
	UserData         map[string]string `json:"user_data"`          // 顶点风格配置
}

// 新增边参数
type EdgeAdd struct {
	Name             string            `json:"name"`               // 边名称
	SourceLabel      string            `json:"source_label"`       // 源顶点类型
	TargetLabel      string            `json:"target_label"`       // 目标顶点类型
	Properties       []string          `json:"properties"`         // 属性列表
	NullableKeys     []string          `json:"nullable_keys"`      // 可空属性列表
	Frequency        string            `json:"frequency"`          // 允许多次连接，可以取值SINGLE和MULTIPLE
	SortKeys         []string          `json:"sort_keys"`          // 当允许关联多次时，指定区分键属性列表
	Ttl              int               `json:"ttl"`                // TTL
	EnableLabelIndex bool              `json:"enable_label_index"` // 启用类型索引,默认为true
	UserData         map[string]string `json:"user_data"`          // 边风格配置
}

// 索引新增参数
type IndexAdd struct {
	Name      string   `json:"name"`       // 索引名称
	BaseType  string   `json:"base_type"`  // 模型类型
	BaseValue string   `json:"base_value"` // 模型名称
	IndexType string   `json:"index_type"` // 索引类型
	Fields    []string `json:"fields"`     // 属性列表
}

// hugegraph-api-append请求接口参数
type PropertiesAppend struct {
	Name         string            `json:"name"`          // 名称
	Properties   []string          `json:"properties"`    // 属性列表
	NullableKeys []string          `json:"nullable_keys"` // 可空属性列表
	UserData     map[string]string `json:"user_data"`     // 风格配置
}
