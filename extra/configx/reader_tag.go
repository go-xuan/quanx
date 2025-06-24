package configx

import (
	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/utils/anyx"
)

// TagReader 标签读取器，获取tag值作为配置值
type TagReader struct {
	Tag string
}

func (r *TagReader) Anchor(tag string) {
	r.Tag = tag
}

func (r *TagReader) Location() string {
	return "tag@" + r.Tag
}

func (r *TagReader) Read(config any) error {
	if err := anyx.SetDefaultValue(config, r.Tag); err != nil {
		return errorx.Wrap(err, "read config from tag error")
	}
	return nil
}
