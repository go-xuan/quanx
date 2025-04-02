package hugegraphx

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/httpx"
	"github.com/go-xuan/quanx/extra/configx"
	"github.com/go-xuan/quanx/extra/nacosx"
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
	FrequencyAll   = "SINGLE,MULTIPLE"                                                                // 单条连接
	DateTypesAll   = "TEXT,BYTE,DATE,INT,LONG,BOOLEAN,DOUBLE,FLOAT,UUID,BLOB,OBJECT,UNKNOWN"          // 属性数据类型
	CardinalityAll = "SINGLE,SET,LIST"                                                                // 属性类型基数
	IdStrategyAll  = "DEFAULT,AUTOMATIC,PRIMARY_KEY,CUSTOMIZE_STRING,CUSTOMIZE_NUMBER,CUSTOMIZE_UUID" // ID策略
	ThicknessAll   = "THICK,NORMAL,FINE"                                                              // 边线条粗细
	SizeAll        = "TINY,SMALL,NORMAL,BIG,HUGE"                                                     // 顶点样式大小
	BaseTypeAll    = "VERTEX_LABEL,EDGE_LABEL"                                                        // 模型类型
	IndexTypeAll   = "SECONDARY,RANGE,SEARCH"                                                         // 索引类型
)

// Config hugegraph配置
type Config struct {
	Host  string `json:"host" yaml:"host" nacos:"hugegraph.host"`    // 主机
	Port  int    `json:"port" yaml:"port" nacos:"hugegraph.port"`    // 端口
	Graph string `json:"graph" yaml:"graph" nacos:"hugegraph.graph"` // 图名称
}

func (c *Config) Format() string {
	return fmt.Sprintf("host=%s port=%d graph=%s", c.Host, c.Port, c.Graph)
}

func (*Config) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FormNacos:
		return &nacosx.Reader{
			DataId: "hugegraph.yaml",
		}
	case configx.FromLocal:
		return &configx.LocalFileReader{
			Name: "hugegraph.yaml",
		}
	default:
		return nil
	}
}

func (c *Config) Execute() error {
	if c.Host != "" && _handler == nil {
		if c.Ping() {
			_handler = &Client{
				config: c,
			}
			log.Info("hugegraph connect success: ", c.Format())
		} else {
			log.Error("hugegraph connect failed: ", c.Format())
		}
	}
	return nil
}

func (c *Config) GremlinUrl() string {
	return fmt.Sprintf("http://%s:%d/gremlin", c.Host, c.Port)
}

func (c *Config) SchemaUrl(uri string) string {
	return fmt.Sprintf("http://%s:%d/graphs/%s/schema/%s", c.Host, c.Port, c.Graph, uri)
}

// Ping gremlin查询API-get请求
func (c *Config) Ping() bool {
	if res, err := httpx.Get(fmt.Sprintf("http://%s:%d/versions", c.Host, c.Port)).Do(); err != nil {
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
