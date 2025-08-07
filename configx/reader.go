package configx

type From string

const (
	FromNacos From = "nacos"
	FromFile  From = "file"
	FromEnv   From = "env"
	FromTag   From = "tag"
)

// CheckReader 检查配置读取器
func CheckReader(reader any, from From) Reader {
	switch from {
	case FromNacos:
		if v, ok := reader.(INacosReader); ok {
			return v.NacosReader()
		}
	case FromFile:
		if v, ok := reader.(IFileReader); ok {
			return v.FileReader()
		}
	case FromEnv:
		if v, ok := reader.(IEnvReader); ok {
			return v.EnvReader()
		}
	case FromTag:
		if v, ok := reader.(ITagReader); ok {
			return v.TagReader()
		}
	}
	return nil
}

// Reader 配置读取器
type Reader interface {
	Anchor(anchor string) // 配置文件锚点
	Location() string     // 配置文件位置
	Read(v any) error     // 配置赋值
}

// INacosReader nacos配置读取器接口
type INacosReader interface {
	NacosReader() Reader
}

// IFileReader 文件配置读取器接口
type IFileReader interface {
	FileReader() Reader
}

// IEnvReader 环境变量配置读取器接口
type IEnvReader interface {
	EnvReader() Reader
}

// ITagReader tag配置读取器接口
type ITagReader interface {
	TagReader() Reader
}
