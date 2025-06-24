package configx

type From string

const (
	FromNacos From = "nacos"
	FromFile  From = "file"
	FromEnv   From = "env"
	FromTag   From = "tag"
)

// Reader 配置读取器
type Reader interface {
	Anchor(anchor string) // 配置文件锚定点
	Location() string     // 配置文件位置
	Read(v any) error     // 配置赋值
}
