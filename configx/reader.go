package configx

import (
	"github.com/go-xuan/utilx/errorx"
)

var (
	defaultFileAnchor = "conf"    // 默认文件读取器锚点
	defaultTagAnchor  = "default" // 默认tag读取器锚点
)

// Reader 配置读取器接口
// 读取器负责从不同来源（如文件、环境变量、远程配置中心等）读取配置数据
type Reader interface {
	// Read 从数据源读取配置并反序列化到目标结构体
	// 根据读取器类型从相应的配置源获取配置数据，并将其解析到传入的结构体中
	Read(v any) error

	// Anchor 设置配置文件锚点
	// 锚点用于标识配置数据在不同配置源中的具体位置
	// 例如：文件读取器使用文件路径作为锚点，而Nacos读取器使用group作为锚点
	Anchor(anchor string)

	// Location 返回配置数据源的位置信息
	// 用于标识配置数据的具体来源位置，便于调试和日志记录
	Location() string
}

// ReaderRead 使用指定的配置读取器读取配置
func ReaderRead(reader Reader, v any) error {
	if reader == nil {
		return errorx.New("reader is nil")
	}
	if err := reader.Read(v); err != nil {
		return errorx.Wrap(err, "reader read error")
	}
	return nil
}
