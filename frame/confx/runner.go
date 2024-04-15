package confx

// 配置器
type Configurator[T any] interface {
	Theme() string   // 配置器主题
	Reader() *Reader // 配置文件读取
	Run() error      // 配置器运行
}

// 配置文件读取
type Reader struct {
	Array       bool   `json:"array" yaml:"array"`             // 是否数组
	FilePath    string `json:"filePath" yaml:"filePath"`       // 本地配置文件路径
	NacosGroup  string `json:"nacosGroup" yaml:"nacosGroup"`   // nacos配置Group
	NacosDataId string `json:"nacosDataId" yaml:"nacosDataId"` // nacos配置ID
	Listen      bool   `json:"listen" yaml:"listen"`           // 是否监听
}
