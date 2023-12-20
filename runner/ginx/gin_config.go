package ginx

// web服务配置
type WebConfig struct {
	Name   string `yaml:"name" default:"app"`       // 服务名
	Host   string `yaml:"host" default:"127.0.0.1"` // 服务host
	Port   int    `yaml:"port" default:"8888"`      // 服务端口
	Prefix string `yaml:"prefix" default:"api"`     // RESTFul api prefix（接口根路由）
	Debug  bool   `yaml:"debug" default:"false"`    // 是否调试环境
}
