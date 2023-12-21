package configx

type Value struct {
	Source string `json:"source" yaml:"source"` // 配置来源
	Where  string `json:"where" yaml:"where"`   // 配置位置
	Name   string `json:"name" yaml:"dataId"`   // 配置文件名
	Type   string `json:"type" yaml:"type"`     // 配置类型
	Listen bool   `json:"listen" yaml:"listen"` // 是否监听
}
