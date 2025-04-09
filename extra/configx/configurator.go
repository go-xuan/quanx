package configx

import "github.com/go-xuan/quanx/base/errorx"

type From string // 配置读取来源

const (
	FromLocal   From = "local"   // 从本地文件读取配置
	FormNacos        = "nacos"   // 从nacos读取配置
	FromDefault      = "default" // 以默认值（tag）生成配置
)

// ReadAndExecute 读取配置文件并运行
func ReadAndExecute(conf Configurator, from From, location ...string) error {
	if reader := conf.Reader(from); reader != nil {
		reader.Location(location...)
		if err := reader.ReadConfig(conf); err != nil {
			return errorx.Wrap(err, "read config failed")
		}
	}
	if err := conf.Execute(); err != nil {
		return errorx.Wrap(err, "configurator execute err")
	}
	return nil
}

// Configurator 配置器
type Configurator interface {
	Format() string     // 配置信息格式化
	Reader(From) Reader // 配置文件读取
	Execute() error     // 配置器运行
}

// Reader 配置读取器
type Reader interface {
	Location(...string) string // 配置文件定位
	ReadConfig(any) error      // 配置读取
}
