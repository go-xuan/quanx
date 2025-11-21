package configx

import (
	"fmt"

	"github.com/go-xuan/utilx/anyx"
	"github.com/go-xuan/utilx/errorx"
	
	"github.com/go-xuan/quanx/constx"
)

// NewTagReader 创建tag读取器
func NewTagReader(tags ...string) *TagReader {
	tag := constx.Default
	if len(tags) > 0 && tags[0] != "" {
		tag = tags[0]
	}
	return &TagReader{
		Tag: tag,
	}
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
	r.Anchor(constx.Default)
	if err := anyx.SetDefaultValue(v, r.Tag); err != nil {
		return errorx.Wrap(err, "read config from tag error")
	}
	return nil
}

func (r *TagReader) Location() string {
	return fmt.Sprintf("tag@%s", r.Tag)
}
