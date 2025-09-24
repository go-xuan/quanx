package quanx

import (
	"github.com/go-xuan/quanx/cachex"
	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/constx"
	"github.com/go-xuan/quanx/gormx"
	"github.com/go-xuan/quanx/logx"
	"github.com/go-xuan/quanx/nacosx"
	"github.com/go-xuan/quanx/redisx"
	"github.com/go-xuan/quanx/serverx"
	"github.com/go-xuan/utilx/anyx"
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
	"github.com/go-xuan/utilx/marshalx"
)

// Config 服务配置
type Config struct {
	Server   *serverx.Config `json:"server" yaml:"server"`     // 服务配置
	Log      *logx.Config    `json:"log" yaml:"log"`           // 日志配置
	Nacos    *nacosx.Config  `json:"nacos" yaml:"nacos"`       // nacos访问配置
	Database *gormx.Configs  `json:"database" yaml:"database"` // 数据源配置
	Redis    *redisx.Configs `json:"redis" yaml:"redis"`       // redis配置
	Cache    *cachex.Configs `json:"cache" yaml:"cache"`       // 缓存配置
}

// Init 初始化配置
func (c *Config) Init(path string) error {
	// 获取预设服务配置
	server := c.PresetServer()
	if !filex.Exists(path) {
		if err := marshalx.Apply(path).Write(path, c); err != nil {
			return errorx.Wrap(err, "write config file error: "+path)
		}
	}

	// 读取配置文件
	if err := marshalx.Apply(path).Read(path, c); err != nil {
		return errorx.Wrap(err, "read config file error: "+path)
	}

	// 初始化nacos
	if err := c.InitNacos(); err != nil {
		return errorx.Wrap(err, "init nacos error")
	}
	// 初始化日志
	if err := c.InitLog(); err != nil {
		return errorx.Wrap(err, "init log error")
	}
	// 初始化数据库
	if err := c.InitDatabase(); err != nil {
		return errorx.Wrap(err, "init database error")
	}
	// 初始化redis
	if err := c.InitRedis(); err != nil {
		return errorx.Wrap(err, "init redis error")
	}
	// 初始化缓存
	if err := c.InitCache(); err != nil {
		return errorx.Wrap(err, "init cache error")
	}

	c.ReloadNacosServer()
	// 加载预设服务配置
	if server != nil {
		c.Server.Port = server.Port
		c.Server.Debug = server.Debug
	}
	return nil
}

// PresetServer 获取预设服务配置
func (c *Config) PresetServer() *serverx.Config {
	if c.Server != nil {
		return &serverx.Config{
			Port:  c.Server.Port,
			Debug: c.Server.Debug,
		}
	}
	return nil
}

// ReloadNacosServer 重新加载nacos服务配置
func (c *Config) ReloadNacosServer() {
	if nacosx.Initialized() {
		var config = &Config{}
		if err := nacosx.NewReader(constx.DefaultConfigFilename).Read(config); err == nil {
			c.Server = config.Server
		}
	}
}

// InitLog 初始化日志配置
func (c *Config) InitLog() error {
	log_ := anyx.IfZero(c.Log, logx.GetConfig())
	if log_.Name == "" || log_.Name == "app" {
		log_.Name = c.Server.Name
	}
	if err := configx.ConfiguratorReadAndExecute(log_); err != nil {
		return errorx.Wrap(err, "run log configurator error")
	}
	c.Log = log_
	return nil
}

// InitNacos 初始化nacos
func (c *Config) InitNacos() error {
	if nacos := c.Nacos; nacos != nil {
		if err := configx.ConfiguratorReadAndExecute(nacos); err != nil {
			return errorx.Wrap(err, "run nacos configurator error")
		}
		// 注册nacos服务实例
		if nacos.EnableNaming() {
			if err := nacosx.RegisterServerInstance(nacos.Group, c.Server); err != nil {
				return errorx.Wrap(err, "nacos register server instance error")
			}
		}
	}
	return nil
}

// InitDatabase 初始化数据库
func (c *Config) InitDatabase() error {
	// 读取数据库配置并初始化
	databases := anyx.IfZero(c.Database, &gormx.Configs{})
	if err := configx.ConfiguratorReadAndExecute(databases); err != nil {
		c.Database = databases
	}
	if !gormx.Initialized() {
		database := &gormx.Config{}
		if err := configx.ConfiguratorReadAndExecute(database); err == nil {
			c.Database = &gormx.Configs{database}
		}
	}
	return nil
}

// InitRedis 初始化缓存
func (c *Config) InitRedis() error {
	// 读取redis配置并初始化
	rds := anyx.IfZero(c.Redis, &redisx.Configs{})
	if err := configx.ConfiguratorReadAndExecute(rds); err == nil {
		c.Redis = rds
	}
	if !redisx.Initialized() {
		redis := &redisx.Config{}
		if err := configx.ConfiguratorReadAndExecute(redis); err == nil {
			c.Redis = &redisx.Configs{redis}
		}
	}
	return nil
}

// InitCache 初始化缓存
func (c *Config) InitCache() error {
	if redisx.Initialized() {
		caches := anyx.IfZero(c.Cache, &cachex.Configs{})
		if err := configx.ConfiguratorReadAndExecute(caches); err == nil {
			c.Cache = caches
		}
		if !cachex.Initialized() {
			var cache = &cachex.Config{}
			if err := configx.ConfiguratorReadAndExecute(cache); err == nil {
				c.Cache = &cachex.Configs{cache}
			}
		}
	}
	return nil
}
