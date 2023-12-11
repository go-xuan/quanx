package hugegraphx

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// hugegraph配置
type Config struct {
	Host  string `json:"host" yaml:"host" nacos:"hugegraph.host"`    // 主机
	Port  int    `json:"port" yaml:"port" nacos:"hugegraph.port"`    // 端口
	Graph string `json:"graph" yaml:"graph" nacos:"hugegraph.graph"` // 图名称
}

// 初始化
func (config *Config) Init() {
	InitHugegraphX(config)
}

// 配置信息格式化
func (config *Config) Format() string {
	return fmt.Sprintf("host=%s port=%d graph=%s", config.Host, config.Port, config.Graph)
}

func (config *Config) GremlinUrl() string {
	return fmt.Sprintf("http://%s:%d/gremlin", config.Host, config.Port)

}

func (config *Config) SchemaUrl() string {
	return fmt.Sprintf("http://%s:%d/graphs/%s/schema/", config.Host, config.Port, config.Graph)
}

// 初始化HugegraphX
func InitHugegraphX(conf *Config) {
	if instance == nil {
		instance = &Handler{
			Config:     conf,
			GremlinUrl: conf.GremlinUrl(),
			SchemaUrl:  conf.SchemaUrl(),
		}
		log.Info("初始化hugegraph连接-成功！", conf.Format())
	}
}
