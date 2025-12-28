package configx

import (
	"fmt"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/reflectx"
)

// SetTagReaderAnchor 设置tag读取器锚点
func SetTagReaderAnchor(anchor string) {
	if anchor != "" {
		defaultTagAnchor = anchor
	}
}

// NewTagReader 创建tag读取器
func NewTagReader() *TagReader {
	return &TagReader{
		Tag: defaultTagAnchor,
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
	r.Anchor(defaultTagAnchor)
	if !reflectx.IsStructPointer(v) {
		return errorx.New("the kind must be struct pointer")
	}
	tag := r.Tag
	val := reflectx.ValueOf(v)
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		if field := val.Field(i); field.IsZero() && field.CanSet() {
			if value := typ.Field(i).Tag.Get(tag); value != "" {
				reflectx.SetValue(field, value)
			}
		}
	}
	return nil
}

func (r *TagReader) Location() string {
	return fmt.Sprintf("tag@%s", r.Tag)
}
