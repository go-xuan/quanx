package configx

import "github.com/go-xuan/quanx/base/errorx"

// Configurator 配置器
type Configurator interface {
	Info() string       // 配置信息格式化
	Reader(From) Reader // 配置文件读取
	Execute() error     // 配置器运行
}

// Reader 配置读取器
type Reader interface {
	Ready(...string)  // 配置器准备
	Check(any) error  // 配置检查
	Read(any) error   // 配置赋值
	Location() string // 配置器准备
}

type From string // 配置读取来源

const (
	FromLocal   From = "local"   // 从本地文件读取配置
	FormNacos        = "nacos"   // 从nacos读取配置
	FromEnv          = "env"     // 从环境变量读取配置
	FromDefault      = "default" // 默认值配置
)

// ReadAndExecute 读取配置文件并运行
func ReadAndExecute(configurator Configurator, from From, location ...string) error {
	if reader := configurator.Reader(from); reader != nil {
		reader.Ready(location...)
		if err := reader.Check(configurator); err != nil {
			return errorx.Wrap(err, "reader check error")
		}
		if err := reader.Read(configurator); err != nil {
			return errorx.Wrap(err, "reader read error")
		}
	}
	if err := configurator.Execute(); err != nil {
		return errorx.Wrap(err, "configurator execute error")
	}
	return nil
}
