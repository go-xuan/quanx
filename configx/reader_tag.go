package configx

import (
	"github.com/go-xuan/utilx/anyx"
	"github.com/go-xuan/utilx/errorx"
)

func NewTagReader(tag ...string) *TagReader {
	if len(tag) > 0 {
		return &TagReader{Tag: tag[0]}
	}
	return &TagReader{}
}

// TagReader tag读取器，获取tag值作为配置值
type TagReader struct {
	Tag string `json:"tag"`
}

func (r *TagReader) Anchor(tag string) {
	if r.Tag == "" {
		r.Tag = tag
	}
}

func (r *TagReader) Read(v any) error {
	r.Anchor("default")
	if err := anyx.SetDefaultValue(v, r.Tag); err != nil {
		return errorx.Wrap(err, "read config from tag error")
	}
	return nil
}

func (r *TagReader) Location() string {
	return "tag@" + r.Tag
}
