package serverx

// Instance 服务实例
type Instance interface {
	GetID() string     // 获取实例ID
	GetName() string   // 获取服务名
	GetIP() string     // 获取IP
	GetPort() int      // 获取端口
	GetDomain() string // 获取域名
	GetStatus() string // 获取状态
}
