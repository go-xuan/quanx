package configx

import (
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/osx"
)

// NewEnvReader 创建环境变量读取器
func NewEnvReader() Reader {
	return &EnvReader{}
}

// EnvReader 环境变量读取器，根据env标签读取环境变量作为配置值
type EnvReader struct{}

func (r *EnvReader) Anchor(string) {
	// 环境变量读取器不需要定义锚点
}

func (r *EnvReader) Read(config any) error {
	if err := osx.SetValueFromEnv(config); err != nil {
		return errorx.Wrap(err, "read env reader failed")
	}
	return nil
}

func (r *EnvReader) Location() string {
	return "env"
}
