package quanx

import (
	"github.com/go-xuan/utilx/anyx"
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"

	"github.com/go-xuan/quanx/cachex"
	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/constx"
	"github.com/go-xuan/quanx/gormx"
	"github.com/go-xuan/quanx/logx"
	"github.com/go-xuan/quanx/nacosx"
	"github.com/go-xuan/quanx/redisx"
	"github.com/go-xuan/quanx/serverx"
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
func (cfg *Config) Init(reader configx.Reader) error {
	var server *serverx.Config
	if s := cfg.Server; s != nil {
		server = &serverx.Config{
			Name:  s.Name,
			Debug: s.Debug,
			Host:  s.Host,
			Http:  s.Http,
			Grpc:  s.Grpc,
		}
	}

	// 初始化配置文件
	if fr, ok := reader.(*configx.FileReader); ok {
		if path := fr.GetPath(); !filex.Exists(path) {
			cfg.Server = serverx.DefaultConfig()
			if err := fr.Write(cfg); err != nil {
				return errorx.Wrap(err, "write config file error")
			}
		}
	}

	// 读取配置文件
	if err := reader.Read(cfg); err != nil {
		return errorx.Wrap(err, "read config error")
	}
	// 初始化nacos
	if err := cfg.initNacos(); err != nil {
		return errorx.Wrap(err, "init nacos error")
	}
	// 注册nacos服务实例
	if err := cfg.registerNacosServer(); err != nil {
		return errorx.Wrap(err, "register server instance error")
	}
	// 初始化日志
	if err := cfg.initLog(); err != nil {
		return errorx.Wrap(err, "init log error")
	}
	// 初始化数据库
	if err := cfg.initDatabase(); err != nil {
		return errorx.Wrap(err, "init database error")
	}
	// 初始化redis
	if err := cfg.initRedis(); err != nil {
		return errorx.Wrap(err, "init redis error")
	}
	// 初始化缓存
	if err := cfg.initCache(); err != nil {
		return errorx.Wrap(err, "init cache error")
	}

	// 覆盖应用配置
	cfg.coverServerConfig(server)

	return nil
}

// 覆盖预设服务配置
func (cfg *Config) coverServerConfig(server *serverx.Config) {
	if server != nil {
		cfg.Server = server
		return
	}
	if nacosx.Initialized() {
		var config = &Config{}
		if err := nacosx.NewReader(constx.DefaultConfigName).Read(config); err == nil {
			cfg.Server = config.Server
			return
		}
	}
}

// 初始化nacos
func (cfg *Config) initNacos() error {
	if nacos := cfg.Nacos; nacos != nil {
		if err := configx.ConfiguratorReadAndExecute(nacos); err != nil {
			return errorx.Wrap(err, "run nacos configurator error")
		}
	}
	return nil
}

// 注册nacos服务实例
func (cfg *Config) registerNacosServer() error {
	if nacos := cfg.Nacos; nacos != nil && nacos.EnableNaming() {
		serverx.Init(&nacosx.ServerCenter{})
		if err := serverx.Register(cfg.Server); err != nil {
			return errorx.Wrap(err, "nacos register server instance error")
		}
	}
	return nil
}

// 初始化日志配置
func (cfg *Config) initLog() error {
	log_ := anyx.IfZero(cfg.Log, logx.GetConfig())
	if log_.Name == "" && cfg.Server != nil {
		log_.Name = cfg.Server.Name
	}
	if err := configx.ConfiguratorReadAndExecute(log_); err != nil {
		return errorx.Wrap(err, "run log configurator error")
	}
	cfg.Log = log_
	return nil
}

// 初始化数据库
func (cfg *Config) initDatabase() error {
	// 读取数据库配置并初始化
	dbs := anyx.IfZero(cfg.Database, &gormx.Configs{})
	if err := configx.ConfiguratorReadAndExecute(dbs); err != nil {
		cfg.Database = dbs
	}
	if !gormx.Initialized() {
		database := &gormx.Config{}
		if err := configx.ConfiguratorReadAndExecute(database); err == nil {
			cfg.Database = &gormx.Configs{database}
		}
	}
	return nil
}

// 初始化缓存
func (cfg *Config) initRedis() error {
	// 读取redis配置并初始化
	rds := anyx.IfZero(cfg.Redis, &redisx.Configs{})
	if err := configx.ConfiguratorReadAndExecute(rds); err == nil {
		cfg.Redis = rds
	}
	if !redisx.Initialized() {
		redis := &redisx.Config{}
		if err := configx.ConfiguratorReadAndExecute(redis); err == nil {
			cfg.Redis = &redisx.Configs{redis}
		}
	}
	return nil
}

// 初始化缓存
func (cfg *Config) initCache() error {
	if redisx.Initialized() {
		caches := anyx.IfZero(cfg.Cache, &cachex.Configs{})
		if err := configx.ConfiguratorReadAndExecute(caches); err == nil {
			cfg.Cache = caches
		}
		if !cachex.Initialized() {
			var cache = &cachex.Config{}
			if err := configx.ConfiguratorReadAndExecute(cache); err == nil {
				cfg.Cache = &cachex.Configs{cache}
			}
		}
	}
	return nil
}
