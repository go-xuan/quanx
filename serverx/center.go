package serverx

var _center Center

// Init 初始化服务中心
func Init(center Center) {
	if _center == nil {
		_center = center
	}
}

// GetCenter 获取服务注册中心
func GetCenter() Center {
	if _center == nil {
		panic("server center not initialized")
	}
	return _center
}

// Center 服务注册中心
type Center interface {
	Register(instance Instance) error           // 注册服务实例
	Deregister(instance Instance) error         // 注销服务实例
	SelectOne(name string) (Instance, error)    // 选择服务单个实例
	SelectList(name string) ([]Instance, error) // 获取服务实例列表
}
