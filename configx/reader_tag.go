package configx

import (
	"github.com/go-xuan/utilx/anyx"
	"github.com/go-xuan/utilx/errorx"
)

// TagReader 标签读取器，获取tag值作为配置值
type TagReader struct {
	Tag string
}

// Anchor 标签读取器锚点为tag名称
func (r *TagReader) Anchor(tag string) {
	if r.Tag == "" {
		r.Tag = tag
	}
}

// Location 配置文件位置
func (r *TagReader) Location() string {
	return "tag@" + r.Tag
}

// Read 从tag中读取值并加载到配置
func (r *TagReader) Read(v any) error {
	if err := anyx.SetDefaultValue(v, r.Tag); err != nil {
		return errorx.Wrap(err, "read config from tag error")
	}
	return nil
}
