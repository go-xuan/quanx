package configx

import "github.com/go-xuan/utilx/errorx"

// Configurator 配置器
type Configurator interface {
	NeedRead() bool     // 是否需要读取
	Reader(From) Reader // 配置读取器
	Execute() error     // 配置器运行
}

// ReadAndExecute 读取配置文件并运行
func ReadAndExecute(configurator Configurator, from From, anchor ...string) error {
	if configurator.NeedRead() {
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
	}
	return nil
}
