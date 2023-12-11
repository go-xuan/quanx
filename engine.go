package quanx

import (
	"errors"
	"fmt"
	"github.com/go-xuan/quanx/utilx/anyx"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-xuan/quanx/gormx"
	"github.com/go-xuan/quanx/logx"
	"github.com/go-xuan/quanx/nacosx"
	"github.com/go-xuan/quanx/redisx"
	"github.com/go-xuan/quanx/utilx/ipx"
	"github.com/go-xuan/quanx/utilx/slicex"
	"github.com/go-xuan/quanx/utilx/structx"
	log "github.com/sirupsen/logrus"
)

var engine *Engine

// 服务启动器
type Engine struct {
	config        *Config                       // 服务依赖配置
	configPath    string                        // 配置文件路径
	initializers  []Initializer                 // 初始化方法
	ginEngine     *gin.Engine                   // gin引擎
	ginLoader     []RouterLoader                // gin路由加载方法
	ginMiddleware []gin.HandlerFunc             // gin中间件
	gormTable     map[string][]gormx.Table[any] // gorm初始化model
}

// 初始化方法
type Initializer func()

// gin路由加载器
type RouterLoader func(router *gin.RouterGroup)

// 应用配置
type Config struct {
	Server   *Server        `yaml:"server"`   // 服务配置
	Log      *logx.Config   `yaml:"log"`      // 日志配置
	Nacos    *nacosx.Config `yaml:"nacos"`    // nacos访问配置
	Database gormx.Configs  `yaml:"database"` // 数据源配置
	Redis    redisx.Configs `yaml:"redis"`    // redis配置
}

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

// 初始化启动器
func GetEngine(load ...RouterLoader) *Engine {
	if engine == nil {
		gin.SetMode(gin.ReleaseMode)
		engine = &Engine{
			config: &Config{
				Server: &Server{},
				Log:    &logx.Config{},
				Nacos:  &nacosx.Config{},
			},
			configPath: "config.yaml",
			ginEngine:  gin.New(),
			gormTable:  make(map[string][]gormx.Table[any]),
		}
		if len(load) > 0 {
			engine.AddGinRouter(load...)
		}
	}
	return engine
}

// 获取服务配置
func GetServerConfig() *Server {
	return GetEngine().config.Server
}

// 服务启动
func (e *Engine) RUN() {
	defer keepAlive()
	// 加载配置
	e.loadConfig()
	// 初始化：日志/nacos/gorm/redis
	e.init()
	// 执行初始化方法
	e.execInitializers()
	// 启动gin
	e.startGin()
}

// 加载服务配置
func (e *Engine) loadConfig() {
	var config = &Config{}
	if err := structx.ReadFileToPointer(config, e.configPath); err != nil {
		log.Error("加载服务配置失败！")
		panic(err)
	}
	if ipx.GetWLANIP() != "" {
		config.Server.Host = ipx.GetWLANIP()
	}
	e.config = config
}

// 初始化日志/nacos/gorm/redis
func (e *Engine) init() {
	var appName = e.config.Server.Name
	// 初始化日志
	var logX X[any] = anyx.IfZero(e.config.Log, &logx.Config{AppName: appName})
	logX.Init()

	// 初始化Nacos
	if e.config.Nacos != nil {
		var nacosX X[any] = e.config.Nacos
		nacosX.Init()
		// 加载Nacos配置
		nacosx.LoadNacosConfig(appName, e.config.Nacos.Config, e.config)
		// 注册Nacos服务
		nacosx.RegisterInstance(
			nacosx.ServerInstance{
				Name:  appName,
				Host:  e.config.Server.Host,
				Port:  e.config.Server.Port,
				Group: e.config.Nacos.NameSpace,
			},
		)
	}

	// 初始化Gorm
	if e.config.Database != nil {
		var gormX X[any] = e.config.Database
		gormX.Init()
		for _, item := range e.config.Database {
			if models, ok := e.gormTable[item.Source]; ok {
				if err := gormx.This().InitGormTable(item.Source, models...); err != nil {
					panic(err)
				}
			}
		}
	}

	// 初始化Redis
	if e.config.Redis != nil {
		var redisX X[any] = e.config.Redis
		redisX.Init()
	}
}

// 执行初始化方法
func (e *Engine) execInitializers() {
	if e.initializers != nil && len(e.initializers) > 0 {
		for _, engineFunc := range e.initializers {
			engineFunc()
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
	e.ginEngine.Use(logx.LoggerToFile(), gin.Recovery())
	if e.ginMiddleware != nil && len(e.ginMiddleware) > 0 {
		e.ginEngine.Use(e.ginMiddleware...)
	}
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

// 设置配置文件
func (e *Engine) SetEngineConfig(path string) {
	e.configPath = path
}

// 添加需要初始化的gormx.Table模型
func (e *Engine) AddTable(dst ...gormx.Table[any]) {
	e.AddTableInSource("default", dst...)
}

// 添加需要某个数据源的gormx.Table模型
func (e *Engine) AddTableInSource(source string, dst ...gormx.Table[any]) {
	e.gormTable[source] = append(e.gormTable[source], dst...)
}

// 添加gin中间件
func (e *Engine) AddGinMiddleware(middleware ...gin.HandlerFunc) {
	e.ginMiddleware = append(e.ginMiddleware, middleware...)
}

// 添加gin的路由加载函数
func (e *Engine) AddGinRouter(load ...RouterLoader) {
	e.ginLoader = append(e.ginLoader, load...)
}

// 添加初始化方法
func (e *Engine) AddInitializer(init ...Initializer) {
	e.initializers = append(e.initializers, init...)
}

// 加载Nacos配置项
func (e *Engine) LoadNacosConfig(config interface{}, dataId string) {
	e.AddInitializer(func() {
		if err := e.loadNacosConfig(config, dataId); err != nil {
			panic(err)
		}
	})
}

// 加载本地配置项
func (e *Engine) LoadLocalConfig(config interface{}, file string) {
	e.AddInitializer(func() {
		if err := e.loadLocalConfig(config, file); err != nil {
			panic(err)
		}
	})
}

// 服务保活
func keepAlive() {
	if err := recover(); err != nil {
		log.Errorf("服务运行失败！错误为 : %s", err)
		return
	}
	select {}
}

// 执行gin的路由加载函数
func (e *Engine) initGinRouter(group *gin.RouterGroup) {
	if e.ginLoader != nil && len(e.ginLoader) > 0 {
		for _, loader := range e.ginLoader {
			loader(group)
		}
	} else {
		log.Warn("ginLoader is empty !")
	}
}

// 加载Nacos配置项
func (e *Engine) loadNacosConfig(config interface{}, dataId string) error {
	if nacosx.This().ConfigClient == nil {
		return errors.New("未初始化nacos配置中心客户端")
	}
	var load = e.config.Nacos.Config
	if load == nil || load.Custom == "" {
		return errors.New("nacos自定义配置项为空，请检查nacos配置信息")
	}
	if !slicex.Contains(slicex.Split(load.Custom, ","), dataId) {
		return errors.New(dataId + "未在nacos中进行配置，请检查nacos配置信息")
	}
	// 获取监听配置项
	var listenMap = make(map[string]bool)
	if load.Listen != "" {
		for _, id := range slicex.Split(load.Listen, ",") {
			listenMap[id] = true
		}
	}
	var item = nacosx.Item{
		Group:  e.config.Server.Name,
		DataId: dataId,
		Listen: listenMap[dataId],
	}
	// 加载微服务配置
	if err := item.LoadConfig(config); err != nil {
		return errors.New("加载nacos自定义配置失败" + err.Error())
	}
	return nil
}

// 加载本地配置项
func (e *Engine) loadLocalConfig(config interface{}, file string) error {
	if err := structx.ReadFileToPointer(config, file); err != nil {
		return err
	}
	return nil
}
