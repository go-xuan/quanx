package configx

import "github.com/go-xuan/quanx/base/errorx"

// Configurator 配置器
type Configurator interface {
	Info() string       // 配置信息
	Execute() error     // 配置器运行
	Reader(From) Reader // 配置读取器
}

// ReadAndExecute 读取配置文件并运行
func ReadAndExecute(configurator Configurator, from From, anchor ...string) error {
	if reader := configurator.Reader(from); reader != nil {
		if len(anchor) > 0 && anchor[0] != "" {
			reader.Anchor(anchor[0])
		}
		if err := reader.Read(configurator); err != nil {
			return errorx.Wrap(err, "reader read error")
		}
		if err := configurator.Execute(); err != nil {
			return errorx.Wrap(err, "configurator execute error")
		}
	}
	return nil
}
