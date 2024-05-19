package hugegraphx

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/core/confx"
	"github.com/go-xuan/quanx/net/httpx"
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

// hugegraph配置
type Hugegraph struct {
	Host  string `json:"host" yaml:"host" nacos:"hugegraph.host"`    // 主机
	Port  int    `json:"port" yaml:"port" nacos:"hugegraph.port"`    // 端口
	Graph string `json:"graph" yaml:"graph" nacos:"hugegraph.graph"` // 图名称
}

// 配置信息格式化
func (h *Hugegraph) ToString() string {
	return fmt.Sprintf("host=%s port=%d graph=%s", h.Host, h.Port, h.Graph)
}

// 配置器标题
func (*Hugegraph) Title() string {
	return "Hugegraph"
}

// 配置文件读取
func (*Hugegraph) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "hugegraph.yaml",
		NacosDataId: "hugegraph.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (h *Hugegraph) Run() error {
	if h.Host != "" && handler == nil {
		var toString = h.ToString()
		if h.Ping() {
			handler = &Handler{Config: h, GremlinUrl: h.GremlinUrl(), SchemaUrl: h.SchemaUrl()}
			log.Info("Hugegraph Connect Successful: ", toString)
		} else {
			log.Error("Hugegraph Connect Failed: ", toString)
		}
	}
	return nil
}

func (h *Hugegraph) GremlinUrl() string {
	return fmt.Sprintf("http://%s:%d/gremlin", h.Host, h.Port)
}

func (h *Hugegraph) SchemaUrl() string {
	return fmt.Sprintf("http://%s:%d/graphs/%s/schema/", h.Host, h.Port, h.Graph)
}

// gremlin查询API-get请求
func (h *Hugegraph) Ping() bool {
	if bytes, err := httpx.Get().Url(fmt.Sprintf("http://%s:%d/versions", h.Host, h.Port)).Do(); err != nil {
		return false
	} else {
		var resp = &PingResp{}
		if err = json.Unmarshal(bytes, &resp); err != nil {
			return false
		}
		if resp.Versions != nil && resp.Versions.Version != "" {
			return true
		}
		return false
	}
}
