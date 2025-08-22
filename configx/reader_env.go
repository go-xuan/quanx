package configx

import (
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/osx"
)

// EnvReader 环境变量读取器，根据env标签读取环境变量作为配置值
type EnvReader struct{}

// Anchor 环境变量读取器锚点暂不定义
func (r *EnvReader) Anchor(string) {}

// Location 配置文件位置
func (r *EnvReader) Location() string {
	return "env"
}

// Read 从环境变量中读取值并加载到配置
func (r *EnvReader) Read(config any) error {
	if err := osx.SetValueFromEnv(config); err != nil {
		return errorx.Wrap(err, "read config from env error")
	}
	return nil
}
