package quanx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/runner/gormx"
	"github.com/go-xuan/quanx/runner/logx"
	"github.com/go-xuan/quanx/runner/nacosx"
	"github.com/go-xuan/quanx/runner/redisx"
	"github.com/go-xuan/quanx/utilx/filex"
	"github.com/go-xuan/quanx/utilx/ipx"
	"github.com/go-xuan/quanx/utilx/structx"
)

// 服务运行器
type Engine struct {
	config        *Config                       // 服务配置
	funcs         []Func                        // 函数运行器
	iRunners      []IRunner[any]                // 运行器
	ginEngine     *gin.Engine                   // gin引擎
	ginLoader     []RouterLoader                // gin路由加载方法
	ginMiddleware []gin.HandlerFunc             // gin中间件
	gormTable     map[string][]gormx.Table[any] // gorm初始化model
	configPath    string                        // 配置文件路径
	notWeb        bool                          // 非web项目
	multiDatabase bool                          // 开启多数据源
	multiRedis    bool                          // 开启多redis数据库
}

// 服务配置
type Config struct {
	Server        *Server              `yaml:"server"`        // 服务配置
	Log           *logx.Log            `yaml:"log"`           // 日志配置
	Nacos         *nacosx.Nacos        `yaml:"nacos"`         // nacos访问配置
	Database      *gormx.Database      `yaml:"database"`      // 数据源配置
	Redis         *redisx.Redis        `yaml:"redis"`         // redis配置
	MultiDatabase *gormx.MultiDatabase `yaml:"multiDatabase"` // 多数据源配置
	MultiRedis    *redisx.MultiRedis   `yaml:"multiRedis"`    // 多redis配置
}

// 函数运行器
type Func func()

// 运行器
type IRunner[T any] interface {
	Name() string                // 运行器名称
	Run() error                  // 运行器运行
	NacosConfig() *nacosx.Config // nacos配置文件
	LocalConfig() string         // 本地配置文件
}

// gin路由加载器
type RouterLoader func(router *gin.RouterGroup)

// 服务配置
type Server struct {
	Name   string `yaml:"name"`                     // 服务名
	Host   string `yaml:"host" default:"127.0.0.1"` // 服务host
	Port   int    `yaml:"port" default:"8888"`      // 服务端口
	Prefix string `yaml:"prefix" default:"api"`     // RESTFul api prefix（接口根路由）
	Debug  bool   `yaml:"debug"`                    // 是否调试环境
}

// 服务地址
func (s *Server) HttpUrl() string {
	return fmt.Sprintf(`http://%s:%d`, s.Host, s.Port)
}

// 服务地址
func (s *Server) ApiUrl() string {
	return s.HttpUrl() + "/" + strings.TrimPrefix(s.Prefix, "/")
}

var engine *Engine

// 服务启动
func (e *Engine) RUN() {
	// 加载配置
	e.loadConfig()
	// 初始化基础配置(日志/nacos/gorm/redis)
	e.initBasic()
	// 运行函数运行器
	e.runFuncs()
	// 运行接口运行器
	e.runIRunners()
	// 启动gin
	if !e.notWeb {
		defer PanicRecover()
		e.startGin()
	}
}

// 初始化运行器
func GetEngine(load ...RouterLoader) *Engine {
	engine = &Engine{
		config:        &Config{},
		configPath:    "conf/config.yaml",
		funcs:         make([]Func, 0),
		iRunners:      make([]IRunner[any], 0),
		ginEngine:     gin.New(),
		ginMiddleware: make([]gin.HandlerFunc, 0),
		gormTable:     make(map[string][]gormx.Table[any]),
	}
	gin.SetMode(gin.ReleaseMode)
	engine.AddGinRouter(load...)
	return engine
}

// 加载服务配置
func (e *Engine) loadConfig() {
	var config = &Config{
		Server: &Server{},
		Log:    &logx.Log{},
	}
	if e.multiDatabase {
		config.MultiDatabase = &gormx.MultiDatabase{}
	} else {
		config.Database = &gormx.Database{}
	}
	if e.multiRedis {
		config.MultiRedis = &redisx.MultiRedis{}
	} else {
		config.Redis = &redisx.Redis{}
	}
	// 读取本地配置
	if filex.Exists(e.configPath) {
		if err := structx.ReadFileToPointer(config, e.configPath); err != nil {
			log.Error("加载服务配置失败！")
			panic(err)
		}
	}
	if ipx.GetWLANIP() != "" {
		config.Server.Host = ipx.GetWLANIP()
	}
	// 初始化nacos
	if config.Nacos != nil {
		e.runIRunner(config.Nacos)
		if !e.notWeb && config.Nacos.OpenNaming() {
			// 注册nacos服务Nacos
			nacosx.RegisterInstance(
				nacosx.ServerInstance{
					Name:  config.Server.Name,
					Host:  config.Server.Host,
					Port:  config.Server.Port,
					Group: config.Nacos.NameSpace,
				},
			)
		}
	}
	e.config = config
}

// 初始化日志/nacos/gorm/redis
func (e *Engine) initBasic() {
	var serverName = e.config.Server.Name
	// 初始化日志
	if e.config.Log != nil {
		e.runIRunner(e.config.Log)
	} else {
		e.runIRunner(&logx.Log{FileName: serverName + ".log"})
	}

	// 连接数据库
	if e.multiDatabase && e.config.MultiDatabase != nil {
		e.runIRunner(e.config.MultiDatabase)
		if gormx.Initialized() && gormx.This().Multi {
			for _, item := range *e.config.MultiDatabase {
				if models, ok := e.gormTable[item.Source]; ok {
					if err := gormx.This().InitGormTable(item.Source, models...); err != nil {
						log.Error("初始化数据库表失败！")
						panic(err)
					}
				}
			}
		}
	}
	if !gormx.Initialized() && e.config.Database != nil {
		e.runIRunner(e.config.Database)
		if models, ok := e.gormTable["default"]; ok {
			if err := gormx.This().InitGormTable("default", models...); err != nil {
				log.Error("初始化数据库表失败！")
				panic(err)
			}
		}
	}
	// 运行redis运行器
	if e.multiRedis && e.config.MultiRedis != nil {
		e.runIRunner(e.config.MultiRedis)
	}
	if !redisx.Initialized() && e.config.Redis != nil {
		e.runIRunner(e.config.Redis)
	}
	return
}

// 运行函数运行器
func (e *Engine) runFuncs() {
	if e.funcs != nil && len(e.funcs) > 0 {
		for _, f := range e.funcs {
			f()
		}
	}
}

// 运行接口运行器
func (e *Engine) runIRunners() {
	if e.iRunners != nil && len(e.iRunners) > 0 {
		for _, runner := range e.iRunners {
			e.runIRunner(runner)
		}
	}
}

// 运行接口运行器
func (e *Engine) runIRunner(runner IRunner[any]) {
	if local := runner.LocalConfig(); local != "" {
		if filex.Exists(local) {
			if err := structx.ReadFileToPointer(runner, local); err != nil {
				log.Info("本地配置文件加载失败! path=", local)
				panic(err)
			}
			log.Info("本地配置文件加载成功! path=", local)
		}
	}
	if config := runner.NacosConfig(); config != nil {
		config.Group = e.config.Server.Name
		if config.Exist() {
			if err := config.LoadConfig(runner); err != nil {
				panic(err)
			}
		}
	}
	if err := runner.Run(); err != nil {
		log.Error(runner.Name(), " -> 运行失败！")
		panic(err)
	}
	log.Info(runner.Name(), " -> 运行完毕！")
}

// 初始化接口运行器
func InitIRunner(runner IRunner[any], nacosGroup string) {
	if local := runner.LocalConfig(); local != "" {
		if filex.Exists(local) {
			if err := structx.ReadFileToPointer(runner, local); err != nil {
				log.Info("本地配置文件加载失败! path=", local)
				panic(err)
			}
			log.Info("本地配置文件加载成功! path=", local)
		}
	}
	if config := runner.NacosConfig(); config != nil {
		config.Group = nacosGroup
		if config.Exist() {
			if err := config.LoadConfig(runner); err != nil {
				panic(err)
			}
		}
	}
}

// 启动gin
func (e *Engine) startGin() {
	if e.config.Server.Debug {
		gin.SetMode(gin.DebugMode)
	}
	if e.ginEngine == nil {
		e.ginEngine = gin.New()
	}
	e.ginEngine.Use(gin.Recovery(), logx.LoggerToFile())
	e.ginEngine.Use(e.ginMiddleware...)
	_ = e.ginEngine.SetTrustedProxies([]string{e.config.Server.Host})
	// 注册根路由，并执行路由注册函数
	var group = e.ginEngine.Group(e.config.Server.Prefix)
	e.initGinRouter(group)
	var port = ":" + strconv.Itoa(e.config.Server.Port)
	log.Info("API接口请求地址: http://" + e.config.Server.Host + port)
	if err := e.ginEngine.Run(port); err != nil {
		log.Error("gin-Engine 运行失败！！！")
		panic(err)
	}
}

func (e *Engine) NotWeb() {
	e.notWeb = true
}

func (e *Engine) MultiDatabase() {
	e.multiDatabase = true
}

func (e *Engine) MultiRedis() {
	e.multiRedis = true
}

// 添加函数运行器
func (e *Engine) AddFStarter(starter ...Func) {
	if len(starter) > 0 {
		e.funcs = append(e.funcs, starter...)
	}
}

// 添加接口运行器
func (e *Engine) AddIRunner(runner ...IRunner[any]) {
	if len(runner) > 0 {
		e.iRunners = append(e.iRunners, runner...)
	}
}

// 设置配置文件
func (e *Engine) SetConfig(path string) {
	e.configPath = path
}

// 添加需要初始化的gormx.Table模型
func (e *Engine) AddTable(dst ...gormx.Table[any]) {
	e.AddTableWithSource("default", dst...)
}

// 添加需要某个数据源的gormx.Table模型
func (e *Engine) AddTableWithSource(source string, dst ...gormx.Table[any]) {
	if len(dst) > 0 {
		e.gormTable[source] = append(e.gormTable[source], dst...)
	}
}

// 添加gin中间件
func (e *Engine) AddGinMiddleware(middleware ...gin.HandlerFunc) {
	if len(middleware) > 0 {
		e.ginMiddleware = append(e.ginMiddleware, middleware...)
	}
}

// 添加gin的路由加载函数
func (e *Engine) AddGinRouter(load ...RouterLoader) {
	if len(load) > 0 {
		e.ginLoader = append(e.ginLoader, load...)
	}
}

// 执行gin的路由加载函数
func (e *Engine) initGinRouter(group *gin.RouterGroup) {
	if e.ginLoader != nil && len(e.ginLoader) > 0 {
		for _, loader := range e.ginLoader {
			loader(group)
		}
	} else {
		log.Warn("engine.ginLoader is empty !")
	}
}

// 加载Nacos配置项（延迟加载，执行RUN()后按添加顺序依次加载）
func (e *Engine) AddNacosConfig(config interface{}, dataId string, listen ...bool) {
	e.AddFStarter(func() {
		if nacosx.This().ConfigClient == nil {
			panic("未初始化nacos配置中心客户端")
		}
		var item = nacosx.Config{Group: e.config.Server.Name, DataId: dataId}
		if len(listen) > 0 {
			item.Listen = listen[0]
		}
		// 加载微服务配置
		if err := item.LoadConfig(config); err != nil {
			panic("加载nacos配置失败" + err.Error())
		}
	})
}

// 加载本地配置项（立即加载）
func (e *Engine) LoadLocalConfig(config interface{}, file string) {
	if err := structx.ReadFileToPointer(config, file); err != nil {
		panic(err)
	}
}

// 服务保活
func PanicRecover() {
	if err := recover(); err != nil {
		log.Errorf("服务运行失败！错误为 : %s", err)
		return
	}
	select {}
}

// 获取服务配置
func GetServerConfig() *Server {
	return GetEngine().config.Server
}
