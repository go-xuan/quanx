package configx

import (
	"github.com/go-xuan/quanx/constx"
	"github.com/go-xuan/utilx/errorx"
)

// DefaultReader 默认配置读取器
func DefaultReader() Reader {
	return &FileReader{
		Dir:  constx.DefaultConfigDir,
		Name: constx.DefaultConfigName,
	}
}

// Reader 配置读取器
type Reader interface {
	Anchor(anchor string) // 配置文件锚点
	Read(v any) error     // 配置读取
	Location() string     // 配置文件位置
}

// ReadWithReader 使用指定的配置读取器读取配置
func ReadWithReader(v any, reader Reader) error {
	if reader == nil {
		return errorx.New("reader is nil")
	}
	if err := reader.Read(v); err != nil {
		return errorx.Wrap(err, "reader read error")
	}
	return nil
}
