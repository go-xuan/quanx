package hugegraphx

import (
	"fmt"
	"strings"

	"github.com/quanxiaoxuan/quanx/common/httpx"
)

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

var CONFIG *Config

// hugegraph配置
type Config struct {
	Host  string `json:"host" yaml:"host" nacos:"hugegraph.host"`    // 主机
	Port  string `json:"port" yaml:"port" nacos:"hugegraph.port"`    // 端口
	Graph string `json:"graph" yaml:"graph" nacos:"hugegraph.graph"` // 图名称
}

// 配置信息格式化
func (conf *Config) Format() string {
	return fmt.Sprintf("host=%s port=%s graph=%s", conf.Host, conf.Port, conf.Graph)
}

func (conf *Config) NewHugegraphCTL() (ctl *Control) {
	ctl = &Control{Config: conf, GremlinUrl: conf.GremlinHttpUrl()}
	return
}

func (conf *Config) GremlinHttpUrl() string {
	httpUrl := strings.Builder{}
	httpUrl.WriteString(httpx.HttpPrefix)
	httpUrl.WriteString(conf.Host)
	httpUrl.WriteString(`:`)
	httpUrl.WriteString(conf.Port)
	httpUrl.WriteString(`/gremlin`)
	return httpUrl.String()
}

func (conf *Config) SchemaHttpUrl(url string) string {
	httpUrl := strings.Builder{}
	httpUrl.WriteString(httpx.HttpPrefix)
	httpUrl.WriteString(conf.Host)
	httpUrl.WriteString(`:`)
	httpUrl.WriteString(conf.Port)
	httpUrl.WriteString(`/graphs/`)
	httpUrl.WriteString(conf.Graph)
	httpUrl.WriteString(`/schema/`)
	httpUrl.WriteString(url)
	return httpUrl.String()
}
