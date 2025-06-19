package configx

import (
	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/osx"
)

// EnvReader 环境变量读取器，根据env标签读取环境变量作为配置值
type EnvReader struct{}

func (r *EnvReader) Anchor(string) {}

func (r *EnvReader) Location() string {
	return "env"
}

func (r *EnvReader) Read(config any) error {
	if err := osx.SetValueFromEnv(config); err != nil {
		return errorx.Wrap(err, "read config from env error")
	}
	return nil
}
