# quanx
基于gin+gorm的web微服务搭建框架
启动项目代码量更少，操作更加简便

## 获取quanx

```
go get github.com/go-xuan/quanx
```

## 服务启动

启动简单，一行代码即可启动一个微服务

```go
package main

import "github.com/go-xuan/quanx"

func main() {
    quanx.NewEngine().RUN()
}
```

### 初始化表结构

```go
func main() {
    engine := quanx.NewEngine()

    // 初始化表结构
    engine.AddTable(
        &User{}, // 需要实现gormx.Tabler接口
    )

    // 启动服务
    engine.RUN()
}

// User 用户表结构必须实现gormx.Tabler接口
type User struct {
    Id           int64     `json:"id" gorm:"type:bigint; not null; comment:用户ID;"`
    Name         string    `json:"name" gorm:"type:varchar(100); not null; comment:姓名;"`
    CreateUserId int64     `json:"createUserId" gorm:"type:bigint; not null; comment:创建人ID;"`
    CreateTime   time.Time `json:"createTime" gorm:"type:timestamp(0); default:now(); comment:创建时间;"`
    UpdateUserId int64     `json:"updateUserId" gorm:"type:bigint; not null; comment:更新人ID;"`
    UpdateTime   time.Time `json:"updateTime" gorm:"type:timestamp(0); default:now(); comment:更新时间;"`
}

// TableName 定义表名
func (User) TableName() string {
    return "user_test"
}
```

### 注册api路由

```go
func main() {
    engine := quanx.NewEngine()

    // 添加gin的路由加载函数
    engine.AddGinRouter(BindApiRouter)

    // 初始化表结构
    engine.AddTable(
        &User{}, // 需要实现 schema.Tabler 接口
    )

    // 启动服务
    engine.RUN()
}

// BindApiRouter 绑定api路由
func BindApiRouter(router *gin.RouterGroup) {
    group := router.Group("/user")
    // 用户表增删改查接口注册，仅一行代码就可以实现CRUD
    ginx.NewCrudApi[User](group, gormx.DB())
}

// User 用户表结构必须实现 schema.Tabler 接口
type User struct {
    Id           int64     `json:"id" gorm:"type:bigint; not null; comment:用户ID;"`
    Name         string    `json:"name" gorm:"type:varchar(100); not null; comment:姓名;"`
    CreateUserId int64     `json:"createUserId" gorm:"type:bigint; not null; comment:创建人ID;"`
    CreateTime   time.Time `json:"createTime" gorm:"type:timestamp(0); default:now(); comment:创建时间;"`
    UpdateUserId int64     `json:"updateUserId" gorm:"type:bigint; not null; comment:更新人ID;"`
    UpdateTime   time.Time `json:"updateTime" gorm:"type:timestamp(0); default:now(); comment:更新时间;"`
}

// TableName 定义表名
func (User) TableName() string {
    return "user_test"
}
```

### 初始化方法

```go
func main() {
    // 初始化服务引擎 
    engine := quanx.NewEngine()

    // 按照添加顺序先后执行 
    engine.AddTaskBefore(Init1, "init1", quanx.TaskRunServer)
    engine.AddTaskAfter(Init2, "init2", "init1")
	
    // 服务启动
    engine.RUN()
}

func Init1() error {
    fmt.Println("执行初始化任务1", time.Now().Format("2006-01-02 15:04:05"))
    return nil
}

func Init2() {
    fmt.Println("执行初始化任务2", time.Now().Format("2006-01-02 15:04:05"))
    return nil
}


```

### 加载自定义配置

```go
func main() {
    // 初始化服务引擎
	engine := quanx.NewEngine()

    // 添加配置器，Config结构体需要实现Configurator接口
    engine.AddConfigurator(Config)

    // 服务启动
    engine.RUN()
}

var Config = &config{}

// 此配置必须实现Configurator配置器接口
type config struct{}

func (c config) Format() string {
    return "配置项格式化文本输出"
}

func (c config) Reader() *configx.Reader {
    return &configx.Reader{
        FilePath:    "config.yaml", // 本地配置文件
        NacosGroup:  "",            // nacos配置分组，默认为服务名
        NacosDataId: "",            // nacos配置ID
        Listen:      false,         // 是否监听
    }
}

func (c config) Execute() error {
    // todo 配置读取后的业务操作
    return nil
}

```

## 服务配置

支持yaml、json、toml、properties等各种配置类型的配置文件。

### 服务配置

quanx框架本身已实现了一些常规配置项的读取和初始化，开发者仅需要在项目代码中添加必要配置文件（默认yaml格式）即可。

#### 主配置

主配置文件路径：conf/config.yaml，此配置必须添加。

```yaml
server:
  name: demo                  # 应用名
  port: 8888                  # 服务端口
  prefix: /demo               # 服务api前缀
```

#### nacos配置

nacos配置文件路径：conf/nacos.yaml，不使用nacos可不添加。

```yaml
address: "127.0.0.1:8848"     # string nacos服务地址,多个以英文逗号分割
username: "nacos"             # string 用户名
password: "nacos"             # string 密码
nameSpace: "demo"             # string 命名空间
mode: 2                       # int 模式（0-仅配置中心；1-仅服务发现；2-配置中心和服务发现）
```

#### 数据库配置

数据库配置文件路径：conf/database.yaml

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
debug: false                  # bool 开启debug模式以及初始化表结构以及数据
```

##### 多数据库

如果想要连接多个数据库，将conf/database.yaml配置文件内容修改为多配置即可

```yaml
- name: default
  enable: true
  type: pgsql
  host: localhost
  port: 5432
  username: postgres
  password: postgres
  database: demo
  debug: false
- name: db1
  enable: true
  type: mysql
  host: localhost
  port: 3306
  username: root
  password: root
  database: demo
  debug: true
```

#### redis配置

redis配置文件路径：conf/redis.yaml

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

如果需要连接多个redis数据源，更新conf/redis.yaml配置文件内容修改为多配置即可

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

每一项配置都需要在代码中使用struct结构体进行声明，并且实现Configurator配置器接口

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

func (d *demo) Info() string {
	return "demo配置"
}

func (d *demo) Reader(from configx.From) confx.Reader {
    switch from {
    case configx.FormNacos: 
        // nacos配置文件读取器
        return &nacosx.Reader{
            DataId: "demo.yaml",
        }
    case configx.FromLocal: 
        // 本地配置文件读取器
        return &configx.LocalReader{
            Name: "demo.yaml",
        }
    default:
        return nil
    }
}

func (d *demo) Execute() error {
	// todo 完成配置读取后需要进行的操作
	fmt.Println(c.Key1)
	fmt.Println(c.Key2)
	fmt.Println(c.Key3)
	return nil
}
```


