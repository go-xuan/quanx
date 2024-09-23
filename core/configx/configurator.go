package configx

import (
	"github.com/go-xuan/quanx/os/errorx"
)

// Configurator 配置器
type Configurator interface {
	ID() string      // 配置器唯一标识
	Format() string  // 配置信息格式化
	Reader() *Reader // 配置文件读取
	Execute() error  // 配置器运行
}

// Reader 配置文件读取器
type Reader struct {
	FilePath    string `json:"filePath" yaml:"filePath"`       // 本地配置文件路径
	NacosDataId string `json:"nacosDataId" yaml:"nacosDataId"` // nacos配置ID
	NacosGroup  string `json:"nacosGroup" yaml:"nacosGroup"`   // nacos配置Group
	Listen      bool   `json:"listen" yaml:"listen"`           // 是否监听
}

// Execute 运行配置器
func Execute(conf Configurator) error {
	if err := conf.Execute(); err != nil {
		return errorx.Wrap(err, "configurator execute err")
	}
	return nil
}
