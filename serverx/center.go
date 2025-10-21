package serverx

var _center Center

// Init 初始化服务中心
func Init(center Center) {
	if _center == nil {
		_center = center
	}
}

// Instance 服务实例
type Instance interface {
	GetName() string // 获取服务名
	GetHost() string // 获取IP
	GetPort() int    // 获取端口
}

// Center 服务注册中心
type Center interface {
	Register(instance Instance) error          // 注册服务实例
	Deregister(instance Instance) error        // 注销服务实例
	SelectOne(name string) (Instance, error)   // 选择单个服务实例
	SelectAll(name string) ([]Instance, error) // 获取全部服务实例
}

func getCenter() Center {
	if _center == nil {
		panic("server center not initialized")
	}
	return _center
}

// Register 注册服务实例
func Register(instance Instance) error {
	return getCenter().Register(instance)
}

// Deregister 注销服务实例
func Deregister(instance Instance) error {
	return getCenter().Deregister(instance)
}

// SelectOne 选择单个服务实例
func SelectOne(name string) (Instance, error) {
	return getCenter().SelectOne(name)
}

// SelectAll 获取全部服务实例
func SelectAll(name string) ([]Instance, error) {
	return getCenter().SelectAll(name)
}
