package hugegraphx

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/httpx"
	"github.com/go-xuan/quanx/extra/configx"
	"github.com/go-xuan/quanx/extra/nacosx"
	"github.com/go-xuan/quanx/types/anyx"
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
	Host  string `json:"host" yaml:"host" default:"localhost"`   // 主机
	Port  int    `json:"port" yaml:"port" default:"8881"`        // 端口
	Graph string `json:"graph" yaml:"graph" default:"hugegraph"` // 图名称
}

func (c *Config) Info() string {
	return fmt.Sprintf("host=%s port=%d graph=%s", c.Host, c.Port, c.Graph)
}

func (*Config) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FormNacos:
		return &nacosx.Reader{
			DataId: "hugegraph.yaml",
		}
	case configx.FromLocal:
		return &configx.LocalReader{
			Name: "hugegraph.yaml",
		}
	default:
		return nil
	}
}

func (c *Config) Execute() error {
	if err := anyx.SetDefaultValue(c); err != nil {
		return errorx.Wrap(err, "set default value error")
	}
	if ok, err := c.Ping(); !ok || err != nil {
		log.WithError(err).Error("hugegraph connect failed: ", c.Info())
		return err
	}
	_client = &Client{
		config:     c,
		gremlinUrl: fmt.Sprintf("http://%s:%d/gremlin", c.Host, c.Port),
		schemaUrl:  fmt.Sprintf("http://%s:%d/graphs/%s/schema/", c.Host, c.Port, c.Graph),
	}
	log.Info("hugegraph connect success: ", c.Info())
	return nil
}

// Ping gremlin查询API-get请求
func (c *Config) Ping() (bool, error) {
	pingUrl := fmt.Sprintf("http://%s:%d/versions", c.Host, c.Port)
	if res, err := httpx.Get(pingUrl).Do(); err != nil || !res.StatusOK() {
		return false, errorx.Wrap(err, "ping request failed")
	} else {
		var resp = struct {
			Versions struct {
				Version string `json:"version"`
				Core    string `json:"core"`
				Gremlin string `json:"gremlin"`
				Api     string `json:"api"`
			} `json:"versions"` // 请求ID
		}{}
		if err = res.Unmarshal(&resp); err != nil {
			return false, errorx.Wrap(err, "resp unmarshal failed")
		}
		log.WithField("version", resp.Versions.Version).
			WithField("core", resp.Versions.Core).
			WithField("gremlin", resp.Versions.Gremlin).
			WithField("api", resp.Versions.Api).
			Info("ping success: ", c.Info())
		return true, nil
	}
}
