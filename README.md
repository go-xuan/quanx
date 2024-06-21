# quanx
基于gin+gorm的web微服务搭建框架
启动项目代码量更少，操作更加简便

## 获取quanx

```
go get github.com/go-xuan/quanx
```

## 服务启动

启动简单，两行代码即可启动一个微服务

```go
package main

import (
	"fmt"
	"time"
	
	"github.com/go-xuan/quanx/app"

	"demo/internal/model/entity"
	"demo/internal/router"
)

func main() {
	// 初始化服务引擎
	var engine = app.NewEngine()

	// 服务启动
	engine.RUN()
}
```

### 初始化表结构

```go
func main() {
	// 初始化服务引擎
	var engine = app.NewEngine()
    
	// 初始化表结构
	engine.AddTable(
		&entity.User{}, // 需要实现gormx.Tabler接口
	)
    
	// 服务启动
	engine.RUN()
}

// gormx.Tabler接口
type Tabler interface {
	TableName() string    // 表名
	TableComment() string // 表注释
	InitData() any        // 表初始数据
}

// 用户表结构必须实现gormx.Tabler接口
type User struct {
	Id           int64     `json:"id" gorm:"type:bigint; not null; comment:用户ID;"`
	Name         string    `json:"name" gorm:"type:varchar; not null; comment:姓名;"`
	CreateUserId int64     `json:"createUserId" gorm:"type:bigint; not null; comment:创建人ID;"`
	CreateTime   time.Time `json:"createTime" gorm:"type:timestamp; default:now(); comment:创建时间;"`
	UpdateUserId int64     `json:"updateUserId" gorm:"type:bigint; not null; comment:更新人ID;"`
	UpdateTime   time.Time `json:"updateTime" gorm:"type:timestamp; default:now(); comment:更新时间;"`
}

// 定义表名
func (User) TableName() string {
	return "t_sys_user"
}

// 定义表注释
func (User) TableComment() string {
	return "用户信息表"
}

// 不为nil时，在初始化表结构之后向表里插入初始化数据
func (User) InitData() any {
	return nil
}
```

### 绑定路由

```go
func main() {
	// 初始化服务引擎
	var engine = app.NewEngine()
    
	// 添加gin的路由加载函数
	engine.AddGinRouter(BindApiRouter)
    
	// 服务启动
	engine.RUN()
}

// 绑定api路由
func BindApiRouter(router *gin.RouterGroup) {
	// 获取默认数据源
	db := gormx.DB()
    
	// 获取具体数据源
	//db := gormx.DB("{{数据源名称}}")

	// 用户信息表-增删改查，一行代码实现CRUD
	ginx.NewCrudApi[entity.User](router.Group("user"), db)
}
```

### 初始化方法

```go
func main() {
	// 初始化服务引擎
	var engine = app.NewEngine(
		app.UseQueue, // 开启任务队列
	)
    
	// 新增初始化方法
	engine.AddCustomFunc(Init1)
	// 或者开启任务队列，使用任务队列
	engine.AddQueueTask("init2", Init2)
    
	// 服务启动
	engine.RUN()
}


func Init1() {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}

func Init2() {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}

```

### 加载自定义配置

```go
func main() {
	// 初始化服务引擎
	var engine = app.NewEngine(
		app.UseQueue, // 开启任务队列
	)

	// 添加配置器，Config结构体需要实现Configurator接口
	engine.AddConfigurator(Config)

	// 服务启动
	engine.RUN()
}

// Configurator配置器接口
type Configurator interface {
	Title() string   // 配置器标题
	Reader() *Reader // 配置文件读取
	Run() error      // 配置器运行
}

var Config *config

// 此配置必须实现Configurator配置器接口
type config struct {
}

func (c config) Title() string {
	return "配置项名称标题"
}

func (c config) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "",    // 本地配置文件
		NacosGroup:  "",    // nacos配置分组，默认为服务名
		NacosDataId: "",    // nacos配置ID
		Listen:      false, // 是否监听
	}
}

func (c config) Run() error {
	// todo 配置读取后的实际操作
	return nil
}
```

## 服务配置

支持yaml、json、toml、properties等各种配置类型的配置文件。

### 服务配置

quanx框架本身已实现了一些常规配置项的读取和初始化，开发者仅需要在项目代码中添加必要配置文件（默认yaml格式）即可。

#### 应用配置

配置文件路径：conf/config.yaml，此配置必须添加。

```yaml
server:
  name: demo                  # 应用名
  port: 8080                  # 服务端口
  prefix: /demo               # 服务api前缀
  debug: true                 # 开启debug模式
```

#### nacos配置

配置文件路径：conf/nacos.yaml，不使用nacos可不添加。

```yaml
address: "127.0.0.1:8848"     # string nacos服务地址,多个以英文逗号分割
username: "nacos"             # string 用户名
password: "nacos"             # string 密码
nameSpace: "demo"             # string 命名空间
mode: 2                       # int 模式（0-仅配置中心；1-仅服务发现；2-配置中心和服务发现）
```

#### 数据库配置

配置文件路径：conf/database.yaml，默认单数据库。

```yaml
source: "default"             # string 数据源名称
enable: false                 # bool 是否启用
type: "mysql"                 # string 数据库类型(mysql/postgres)
host: "127.0.0.1"             # string host
port: 5432                    # int 端口
username: "root"              # string 用户名
password: "root"              # string 密码
database: ""                  # string 数据库名
schema: ""                    # string 模式名（postgres）
debug: false                  # bool 开启debug
init: false                   # bool 是否初始化表结构以及数据
```

##### 多数据源

如果想要连接多个数据库，需要在启动时开启多数据源：

```go
func main() {
	var engine = app.NewEngine(
		app.MultiDatabase, // 开启多数据源
	)
}
```

同时更新conf/database.yaml配置文件内容为：

```yaml
- name: default
  enable: true
  type: 
  host: 
  port: 
  username: 
  password: 
  database: 
  debug: 
- name: db1
  enable: 
  type: 
  host: 
  port: 
  username: 
  password: 
  database: 
  debug: true
......
```

#### redis配置

配置文件路径：conf/redis.yaml，默认单redis数据库。

```yaml
source: "default"             # string 数据源名称
enable: false                 # bool 是否启用
host: "127.0.0.1"             # string host
port: 6379                    # int 端口
password: ""                  # string 密码
database: 0                   # int 数据库
mode: 0                       # int 模式（0-单机；1-集群），默认单机模式
```

##### 多redis源

如果需要连接多个redis数据源，需要在启动时开启多数据源：

```go
func main() {
	var engine = app.NewEngine(
		app.MultiRedis, // 开启多redis数据源
	)
}
```

更新conf/redis.yaml配置文件内容为：

```yaml
- name: default
  enable: 
  host: 
  port: 
  password: 
  database: 
  mode: 0
- name: redis_db1
  enable: 
  host: 
  port: 
  password: 
  database: 
  mode: 0
......
```

#### 自定义配置

每一项配置都需要在go代码中使用struct进行声明，而且结构体并实现Configurator配置器接口。

demo.yaml：

```yaml
key1: 123
key2: "456"
key3:
  - "abc"
  - "def"
```

对应结构体：

```go
// 此配置需要实现Configurator配置器接口
type demo struct {
	Key1 int      `json:"key1" yaml:"key1"`
	Key2 string   `json:"key2" yaml:"key2"`
	Key3 []string `json:"key3" yaml:"key3"`
}

func (d demo) Title() string {
	return "demo配置"
}

func (d demo) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "demo.yaml", // 本地配置文件
		NacosGroup:  "demo",      // nacos配置分组，默认为服务名
		NacosDataId: "demo.yaml", // nacos配置ID
		Listen:      false,       // 是否监听
	}
}

func (d demo) Run() error {
	// todo 配置读取后的实际操作
	fmt.Println(c.Key1)
	fmt.Println(c.Key2)
	fmt.Println(c.Key3)
	return nil
}
```



##### 本地配置

当服务启动时不启用nacos，并且配置项对应结构实现Configurator接口时，Reader()方法返回的Reader.FilePath不为空。

```go
func main() {
	// 初始化服务启动引擎
	// 启动参数不加app.EnableNacos即表示不使用nacos
	var engine = app.NewEngine()
}

func (d demo) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "demo.yaml", // 本地配置文件
	}
}

```

##### Nacos配置

当服务启动启用nacos，并且配置项对应结构实现Configurator接口时，Reader()方法返回的Reader.NacosDataId不为空。

```go
func main() {
	// 初始化服务启动引擎
	var engine = app.NewEngine(
		app.EnableNacos, // 启用nacos
	)
}

func (d demo) Reader() *confx.Reader {
	return &confx.Reader{
		NacosGroup:  "demo",      // nacos配置分组，默认为服务名
		NacosDataId: "demo.yaml", // nacos配置ID
		Listen:      false,       // 是否监听
	}
}
```

