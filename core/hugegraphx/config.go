package hugegraphx

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/core/configx"
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
	CardinalityAll = "SINGLE,Set,LIST"                                                                // 属性类型基数
	IdStrategyAll  = "DEFAULT,AUTOMATIC,PRIMARY_KEY,CUSTOMIZE_STRING,CUSTOMIZE_NUMBER,CUSTOMIZE_UUID" // ID策略
	ThicknessAll   = "THICK,NORMAL,FINE"                                                              // 边线条粗细
	SizeAll        = "TINY,SMALL,NORMAL,BIG,HUGE"                                                     // 顶点样式大小
	BaseTypeAll    = "VERTEX_LABEL,EDGE_LABEL"                                                        // 模型类型
	IndexTypeAll   = "SECONDARY,RANGE,SEARCH"                                                         // 索引类型
)

func NewConfigurator(conf *Config) configx.Configurator {
	return conf
}

// Config hugegraph配置
type Config struct {
	Host  string `json:"host" yaml:"host" nacos:"hugegraph.host"`    // 主机
	Port  int    `json:"port" yaml:"port" nacos:"hugegraph.port"`    // 端口
	Graph string `json:"graph" yaml:"graph" nacos:"hugegraph.graph"` // 图名称
}

func (*Config) ID() string {
	return "hugegraph"
}

func (h *Config) Format() string {
	return fmt.Sprintf("host=%s port=%d graph=%s", h.Host, h.Port, h.Graph)
}

func (*Config) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "hugegraph.yaml",
		NacosDataId: "hugegraph.yaml",
		Listen:      false,
	}
}

func (h *Config) Execute() error {
	if h.Host != "" && _handler == nil {
		if h.Ping() {
			_handler = &Handler{
				config:     h,
				gremlinUrl: h.GremlinUrl(),
				schemaUrl:  h.SchemaUrl()}
			log.Info("hugegraph connect successfully: ", h.Format())
		} else {
			log.Error("hugegraph connect failed: ", h.Format())
		}
	}
	return nil
}

func (h *Config) GremlinUrl() string {
	return fmt.Sprintf("http://%s:%d/gremlin", h.Host, h.Port)
}

func (h *Config) SchemaUrl() string {
	return fmt.Sprintf("http://%s:%d/graphs/%s/schema/", h.Host, h.Port, h.Graph)
}

// Ping gremlin查询API-get请求
func (h *Config) Ping() bool {
	if res, err := httpx.Get(fmt.Sprintf("http://%s:%d/versions", h.Host, h.Port)).Do(); err != nil {
		return false
	} else if res.StatusOK() {
		var resp = &PingResp{}
		if err = res.Unmarshal(resp); err != nil {
			return false
		}
		if resp.Versions != nil && resp.Versions.Version != "" {
			return true
		}
	}
	return false
}
