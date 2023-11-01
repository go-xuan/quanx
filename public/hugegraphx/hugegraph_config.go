package hugegraphx

import (
	"fmt"
	"github.com/quanxiaoxuan/quanx/common/constx"
	"strconv"
	"strings"
)

// hugegraph配置
type Config struct {
	Host  string `json:"host" yaml:"host" nacos:"hugegraph.host"`    // 主机
	Port  int    `json:"port" yaml:"port" nacos:"hugegraph.port"`    // 端口
	Graph string `json:"graph" yaml:"graph" nacos:"hugegraph.graph"` // 图名称
}

// 配置信息格式化
func (conf *Config) Format() string {
	return fmt.Sprintf("host=%s port=%d graph=%s", conf.Host, conf.Port, conf.Graph)
}

func (conf *Config) GremlinHttpUrl() string {
	httpUrl := strings.Builder{}
	httpUrl.WriteString(constx.HttpPrefix)
	httpUrl.WriteString(conf.Host)
	httpUrl.WriteString(`:`)
	httpUrl.WriteString(strconv.Itoa(conf.Port))
	httpUrl.WriteString(`/gremlin`)
	return httpUrl.String()
}

func (conf *Config) SchemaHttpUrl(url string) string {
	httpUrl := strings.Builder{}
	httpUrl.WriteString(constx.HttpPrefix)
	httpUrl.WriteString(conf.Host)
	httpUrl.WriteString(`:`)
	httpUrl.WriteString(strconv.Itoa(conf.Port))
	httpUrl.WriteString(`/graphs/`)
	httpUrl.WriteString(conf.Graph)
	httpUrl.WriteString(`/schema/`)
	httpUrl.WriteString(url)
	return httpUrl.String()
}
