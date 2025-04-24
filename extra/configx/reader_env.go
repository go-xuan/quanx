package configx

import (
	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/osx"
)

type EnvReader struct {
}

func (r *EnvReader) Ready(...string) {
}

func (r *EnvReader) Check(_ any) error {
	return nil
}

func (r *EnvReader) Read(config any) error {
	if err := osx.SetValueFromEnv(config); err != nil {
		return errorx.Wrap(err, "read config from env error")
	}
	return nil
}

func (r *EnvReader) Location() string {
	return "env"
}
