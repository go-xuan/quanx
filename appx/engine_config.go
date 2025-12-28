package appx

import (
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"

	"github.com/go-xuan/quanx/cachex"
	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/gormx"
	"github.com/go-xuan/quanx/logx"
	"github.com/go-xuan/quanx/nacosx"
	"github.com/go-xuan/quanx/redisx"
	"github.com/go-xuan/quanx/serverx"
)

// Config 应用配置
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
	// 预制配置
	var server *serverx.Config
	if srv := cfg.Server; srv != nil {
		server = &serverx.Config{}
		server.Cover(srv)
	}

	// 初始化配置文件
	if fr, ok := reader.(*configx.FileReader); ok {
		if path := fr.GetPath(); !filex.Exists(path) {
			if server != nil {
				cfg.Server = server
			} else {
				cfg.Server = serverx.DefaultConfig()
			}
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
	// 初始化服务配置
	if err := cfg.initServer(server); err != nil {
		return errorx.Wrap(err, "init server error")
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

	return nil
}

// 初始化nacos
func (cfg *Config) initNacos() error {
	if nacos := cfg.Nacos; nacos != nil {
		if err := configx.LoadConfigurator(nacos); err != nil {
			return errorx.Wrap(err, "run nacos configurator error")
		}
	}
	return nil
}

// 初始化服务配置（生效优先级：预制配置 > nacos > 配置文件）
func (cfg *Config) initServer(server *serverx.Config) error {
	// 当前cfg.Server为配置文件
	if cfg.Server != nil {
		// 1.读取nacos配置，有则覆盖
		if nacosx.Initialized() {
			config := &Config{}
			if err := nacosx.NewReader(DefaultConfigName).Read(config); err == nil {
				cfg.Server.Cover(config.Server)
			}
		}

		// 2.覆盖预制配置
		cfg.Server.Cover(server)

		// 如果nacos配置启用了服务发现，则自动注册当前服务
		if cfg.Nacos != nil && cfg.Nacos.EnableNaming() {
			serverx.Init(&nacosx.ServerCenter{})
			if err := serverx.Register(cfg.Server); err != nil {
				return errorx.Wrap(err, "register nacos server instance error")
			}
		}
	}
	return nil
}

// 初始化日志配置
func (cfg *Config) initLog() error {
	log := cfg.Log
	if cfg.Log == nil {
		log = logx.GetConfig()
	}
	if log.Name == "" && cfg.Server != nil {
		log.Name = cfg.Server.Name
	}
	if err := configx.LoadConfigurator(log); err != nil {
		return errorx.Wrap(err, "run log configurator error")
	}
	cfg.Log = log
	return nil
}

// 初始化数据库
func (cfg *Config) initDatabase() error {
	// 读取数据库配置并初始化
	dbs := cfg.Database
	if dbs == nil {
		dbs = &gormx.Configs{}
	}
	if err := configx.LoadConfigurator(dbs); err != nil {
		cfg.Database = dbs
	}
	if !gormx.Initialized() {
		database := &gormx.Config{}
		if err := configx.LoadConfigurator(database); err == nil {
			cfg.Database = &gormx.Configs{database}
		}
	}
	return nil
}

// 初始化redis
func (cfg *Config) initRedis() error {
	// 读取redis配置并初始化
	rds := cfg.Redis
	if rds == nil {
		rds = &redisx.Configs{}
	}
	if err := configx.LoadConfigurator(rds); err == nil {
		cfg.Redis = rds
	}
	if !redisx.Initialized() {
		redis := &redisx.Config{}
		if err := configx.LoadConfigurator(redis); err == nil {
			cfg.Redis = &redisx.Configs{redis}
		}
	}
	return nil
}

// 初始化缓存
func (cfg *Config) initCache() error {
	if redisx.Initialized() {
		caches := cfg.Cache
		if caches == nil {
			caches = &cachex.Configs{}
		}
		if err := configx.LoadConfigurator(caches); err == nil {
			cfg.Cache = caches
		}
		if !cachex.Initialized() {
			var cache = &cachex.Config{}
			if err := configx.LoadConfigurator(cache); err == nil {
				cfg.Cache = &cachex.Configs{cache}
			}
		}
	}
	return nil
}
