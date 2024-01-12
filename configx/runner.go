package configx

// 运行器
type IRunner[T any] interface {
	Run() error            // 运行器运行
	Title() string         // 运行器标题
	ConfigReader() *Reader // 配置文件读取
}

// 配置文件读取
type Reader struct {
	FilePath    string `json:"filePath" yaml:"filePath"`       // 本地配置文件路径
	NacosGroup  string `json:"nacosGroup" yaml:"nacosGroup"`   // nacos配置Group
	NacosDataId string `json:"nacosDataId" yaml:"nacosDataId"` // nacos配置ID
	Listen      bool   `json:"listen" yaml:"listen"`           // 是否监听
}
