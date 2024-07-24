package confx

import log "github.com/sirupsen/logrus"

// Configurator 配置器
type Configurator interface {
	Title() string   // 配置器标题
	Reader() *Reader // 配置文件读取
	Run() error      // 配置器运行
}

// Reader 配置文件读取器
type Reader struct {
	FilePath    string `json:"filePath" yaml:"filePath"`       // 本地配置文件路径
	NacosGroup  string `json:"nacosGroup" yaml:"nacosGroup"`   // nacos配置Group
	NacosDataId string `json:"nacosDataId" yaml:"nacosDataId"` // nacos配置ID
	Listen      bool   `json:"listen" yaml:"listen"`           // 是否监听
}

// RunConfigurator 运行配置器
func RunConfigurator(conf Configurator) {
	if err := conf.Run(); err != nil {
		log.Error(conf.Title(), " Run Failed!")
		panic(err)
	}
	log.Info(conf.Title(), " Run Completed!")
}
