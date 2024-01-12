package quanx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xuan/quanx/utilx/stringx"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strconv"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/importx/gormx"
	"github.com/go-xuan/quanx/importx/logx"
	"github.com/go-xuan/quanx/importx/nacosx"
	"github.com/go-xuan/quanx/importx/redisx"
	"github.com/go-xuan/quanx/utilx/anyx"
	"github.com/go-xuan/quanx/utilx/ipx"
	"github.com/go-xuan/quanx/utilx/marshalx"
)

var engine *Engine

// 初始化Engine
func GetEngine(modes ...Flag) *Engine {
	if engine == nil {
		engine = &Engine{
			config:         &Config{},
			customFuncs:    make([]CustomFunc, 0),
			iRunners:       make([]configx.IRunner[any], 0),
			ginMiddlewares: make([]gin.HandlerFunc, 0),
			gormTables:     make(map[string][]gormx.Table[any]),
			flag:           make(map[Flag]bool),
		}
		gin.SetMode(gin.ReleaseMode)
		engine.SetMode(modes...)
	}
	return engine
}

// 获取服务配置
func GetServer() *Server {
	return GetEngine().config.server
}

// 服务运行器
type Engine struct {
	flag           map[Flag]bool                 // 服务运行标识
	config         *Config                       // 服务配置 使用 loadConfig() 将配置文件加载到此
	configDir      string                        // 服务配置文件文件夹, 使用 SetConfigDir() 设置配置文件路径
	iRunners       []configx.IRunner[any]        // 自定义运行器，使用 AddIRunner() 添加运行器对象，被添加对象必须实现 configx.IRunner 接口且为指针类型
	customFuncs    []CustomFunc                  // 自定义初始化函数 使用 AddCustomFunc() 添加自定义函数
	ginEngine      *gin.Engine                   // gin框架实例
	ginLoaders     []RouterLoader                // gin路由的预加载方法，使用 AddGinRouter() 添加自行实现的路由注册方法
	ginMiddlewares []gin.HandlerFunc             // gin中间件的预加载方法，使用 AddGinRouter() 添加gin中间件
	gormTables     map[string][]gormx.Table[any] // gorm表结构实体对象，使用 AddTable() / AddSourceTable() 添加表结构初始化任务列表，需要实现 gormx.Table 接口
}

// 服务配置
type Config struct {
	server     *Server              `yaml:"server"`     // 服务配置
	log        *logx.Log            `yaml:"log"`        // 日志配置
	nacos      *nacosx.Nacos        `yaml:"nacos"`      // nacos访问配置
	database   *gormx.Database      `yaml:"database"`   // 数据源配置
	redis      *redisx.Redis        `yaml:"redis"`      // redis配置
	multiDB    *gormx.MultiDatabase `yaml:"multiDB"`    // 多数据源配置
	multiRedis *redisx.MultiRedis   `yaml:"multiRedis"` // 多redis配置
}

// gin路由加载器
type RouterLoader func(router *gin.RouterGroup)

// 自定义函数
type CustomFunc func()

// 服务启动标识
type Flag int

const (
	DisableGin    Flag = iota // 关闭gin
	EnableNacos               // 启用nacos
	MultiDB                   // 启用多数据源
	MultiRedis                // 启用多Redis
	IfLoadConfig              // 是否已加载配置
	IfInitBasic               // 是否已初始化基础配置
	IfRunFunc                 // 是否已执行自定义函数
	IfRunIRunners             // 是否已执行接口运行器
)

// 服务配置
type Server struct {
	Name   string `yaml:"name"`                     // 服务名
	Host   string `yaml:"host" default:"127.0.0.1"` // 服务host
	Port   int    `yaml:"port" default:"8888"`      // 服务端口
	Prefix string `yaml:"prefix" default:"app"`     // RESTFul api prefix（接口根路由）
	Debug  bool   `yaml:"debug"`                    // 是否调试环境
}

// 服务地址
func (s *Server) HttpUrl() string {
	return fmt.Sprintf(`http://%s:%d/%s`, s.Host, s.Port, stringx.TrimPrefix(s.Prefix, "/"))
}

// 服务运行
func (e *Engine) RUN() {
	e.PREPARE()  // 服务准备
	e.STARTGIN() // 服务启动
}

// 服务准备
func (e *Engine) PREPARE() {
	// 加载配置
	if !e.flag[IfLoadConfig] {
		e.loadConfig()
	}
	// 初始化基础配置(日志/nacos/gorm/redis)
	if !e.flag[IfInitBasic] {
		e.initBasic()
	}
	// 执行接口运行器
	if !e.flag[IfRunIRunners] {
		e.runIRunners()
	}
	// 执行自定义函数
	if !e.flag[IfRunFunc] {
		e.runCustomFunc()
	}
}

// 服务启动
func (e *Engine) STARTGIN() {
	if !e.flag[DisableGin] {
		defer PanicRecover()
		e.startGin()
	}
}

// 加载服务配置
func (e *Engine) loadConfig() {
	var config = &Config{server: &Server{}}
	// 读取本地配置
	var path = e.GetConfigPath("config.yaml")
	if err := marshalx.LoadFromFile(config, path); err != nil {
		log.Error("加载服务配置失败!")
		panic(err)
	}
	if ipx.GetWLANIP() != "" {
		config.server.Host = ipx.GetWLANIP()
	}
	// 初始化nacos
	if e.flag[EnableNacos] && config.nacos != nil {
		e.runIRunner(config.nacos, true)
		if config.nacos.EnableNaming() {
			// 注册nacos服务Nacos
			nacosx.RegisterInstance(
				nacosx.ServerInstance{
					Name:  config.server.Name,
					Host:  config.server.Host,
					Port:  config.server.Port,
					Group: config.nacos.NameSpace,
				},
			)
		}
	}
	e.config = config
	e.flag[IfLoadConfig] = true
}

// 初始化日志/nacos/gorm/redis
func (e *Engine) initBasic() {
	var serverName = e.config.server.Name
	// 初始化日志
	e.runIRunner(anyx.IfZero(e.config.log, logx.New(serverName)), true)
	// 连接数据库
	if e.flag[MultiDB] {
		e.config.multiDB = &gormx.MultiDatabase{}
		e.runIRunner(e.config.multiDB)
		if gormx.Initialized() && gormx.This().Multi {
			for _, item := range *e.config.multiDB {
				if dst, ok := e.gormTables[item.Source]; ok {
					if err := gormx.This().InitGormTable(item.Source, dst...); err != nil {
						log.Error("初始化数据库表失败!")
						panic(err)
					}
				}
			}
		}
	} else if !gormx.Initialized() {
		e.config.database = &gormx.Database{}
		e.runIRunner(e.config.database)
		if dst, ok := e.gormTables["default"]; ok {
			if err := gormx.This().InitGormTable("default", dst...); err != nil {
				log.Error("初始化数据库表失败!")
				panic(err)
			}
		}
	}
	// 连接redis
	if e.flag[MultiRedis] {
		e.config.multiRedis = &redisx.MultiRedis{}
		e.runIRunner(e.config.multiRedis)
	} else if !redisx.Initialized() {
		e.config.redis = &redisx.Redis{}
		e.runIRunner(e.config.redis)
	}
	e.flag[IfInitBasic] = true
}

// 添加自定义函数
func (e *Engine) AddCustomFunc(starter ...CustomFunc) {
	if len(starter) > 0 {
		e.customFuncs = append(e.customFuncs, starter...)
	}
}

// 执行自定义函数
func (e *Engine) runCustomFunc() {
	if e.customFuncs != nil && len(e.customFuncs) > 0 {
		for _, f := range e.customFuncs {
			f()
		}
	}
	e.flag[IfRunFunc] = true
}

// 添加接口运行器
func (e *Engine) AddIRunner(runner ...configx.IRunner[any]) {
	if len(runner) > 0 {
		e.iRunners = append(e.iRunners, runner...)
	}
}

// 运行接口运行器
func (e *Engine) runIRunners() {
	if e.iRunners != nil && len(e.iRunners) > 0 {
		for _, runner := range e.iRunners {
			e.runIRunner(runner)
		}
	}
	e.flag[IfRunIRunners] = true
}

// 运行接口运行器
func (e *Engine) runIRunner(runner configx.IRunner[any], must ...bool) {
	var ok bool
	if len(must) > 0 {
		ok = must[0]
	}
	if reader := runner.ConfigReader(); reader != nil {
		if e.flag[EnableNacos] {
			reader.NacosGroup = e.config.server.Name
			if err := nacosx.NewConfig(reader.NacosGroup, reader.NacosDataId).LoadConfig(runner); err == nil {
				ok = true
			}
		} else {
			if err := marshalx.LoadFromFile(runner, e.GetConfigPath(reader.FilePath)); err == nil {
				ok = true
			}
		}
	}
	if ok {
		if err := runner.Run(); err != nil {
			log.Error(runner.Title() + " error!")
			panic(err)
		}
		log.Info(runner.Title() + " completed!")
	}
}

// 启动gin
func (e *Engine) startGin() {
	if e.config.server.Debug {
		gin.SetMode(gin.DebugMode)
	}
	if e.ginEngine == nil {
		e.ginEngine = gin.New()
	}
	e.ginEngine.Use(gin.Recovery(), logx.LoggerToFile())
	e.ginEngine.Use(e.ginMiddlewares...)
	_ = e.ginEngine.SetTrustedProxies([]string{e.config.server.Host})
	// 注册服务根路由，并执行路由注册函数
	var group = e.ginEngine.Group(e.config.server.Prefix)
	e.InitGinLoader(group)
	var port = ":" + strconv.Itoa(e.config.server.Port)
	log.Info("API接口请求地址: http://" + e.config.server.Host + port)
	if err := e.ginEngine.Run(port); err != nil {
		log.Error("gin-Engine 运行失败!!!")
		panic(err)
	}
}

// 设置模式
func (e *Engine) SetMode(flags ...Flag) {
	if len(flags) > 0 {
		for _, flag := range flags {
			e.flag[flag] = true
		}
	}
}

// 设置配置文件
func (e *Engine) SetConfigDir(dir string) {
	e.configDir = dir
}

// 设置配置文件
func (e *Engine) GetConfigPath(name string) string {
	return anyx.IfZeroElse(e.configDir, name, filepath.Join(e.configDir, name))
}

// 添加需要初始化的 gormx.Table 模型
func (e *Engine) AddTable(dst ...gormx.Table[any]) {
	e.AddSourceTable("default", dst...)
}

// 添加需要某个数据源的gormx.Table模型
func (e *Engine) AddSourceTable(source string, dst ...gormx.Table[any]) {
	if len(dst) > 0 {
		e.gormTables[source] = append(e.gormTables[source], dst...)
	}
}

// 添加gin中间件
func (e *Engine) AddGinMiddleware(middleware ...gin.HandlerFunc) {
	if len(middleware) > 0 {
		e.ginMiddlewares = append(e.ginMiddlewares, middleware...)
	}
}

// 添加gin的路由加载函数
func (e *Engine) AddGinRouter(load ...RouterLoader) {
	if len(load) > 0 {
		e.ginLoaders = append(e.ginLoaders, load...)
	}
}

// 执行gin的路由加载函数
func (e *Engine) InitGinLoader(group *gin.RouterGroup) {
	if e.ginLoaders != nil && len(e.ginLoaders) > 0 {
		for _, loader := range e.ginLoaders {
			loader(group)
		}
	} else {
		log.Warn("engine.ginLoaders is empty !")
	}
}

// 初始化本地配置项（立即加载）
func (e *Engine) InitLocalConfig(config interface{}, filePath string) {
	if err := marshalx.LoadFromFile(config, filePath); err != nil {
		panic(err)
	}
}

// 初始化Nacos配置项（以自定义函数的形式延迟加载）
func (e *Engine) InitNacosConfig(config interface{}, dataId string, listen ...bool) {
	e.AddCustomFunc(func() {
		if nacosx.This().ConfigClient == nil {
			panic("未初始化nacos配置中心客户端")
		}
		var item = nacosx.NewConfig(e.config.server.Name, dataId)
		if len(listen) > 0 {
			item.Listen = listen[0]
		}
		// 加载微服务配置
		if err := item.LoadConfig(config); err != nil {
			panic("加载nacos配置失败" + err.Error())
		}
	})
}

// 服务保活
func PanicRecover() {
	if err := recover(); err != nil {
		log.Errorf("服务运行失败!错误为 : %s", err)
		return
	}
	select {}
}
