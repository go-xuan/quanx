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

import (
	"context"

	"github.com/go-xuan/quanx/appx"
)

func main() {
	appx.NewEngine().RUN(context.Background())
}

```

### 初始化表结构

```go
package main

import (
	"context"
	"time"

	"github.com/go-xuan/quanx/appx"
)

func main() {
	appx.NewEngine(
		appx.AddTable(&User{}), // 初始化表结构
	).RUN(context.Background())
}

// User 用户表结构必须实现 dbx.Tabler 接口
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
package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-xuan/quanx/appx"
	"github.com/go-xuan/quanx/ginx"
	"github.com/go-xuan/quanx/serverx"
)

func main() {
	appx.NewEngine(
		appx.AddServer(HttpServer()), // 添加http服务
		appx.AddTable(&User{}),       // 初始化表结构
	).RUN(context.Background())
}

// HttpServer 创建http服务
func HttpServer() *serverx.HttpServer {
	return serverx.NewHttpServer(ginx.NewHttpServer(
		ginx.SetDebugMode, // 开启调试模式
		BindRouter,        // 绑定路由
	))
}

// BindRouter 绑定api路由
func BindRouter(engine *gin.Engine) {
	group := engine.Group("/user")
	// 用户表增删改查接口注册，仅一行代码就可以实现CRUD
	ginx.NewCrudApi[User](group, "default")
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


### 加载自定义配置

```go
package main

import (
	"context"

	"github.com/go-xuan/quanx/appx"
	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

func main() {
	// 初始化服务引擎
	appx.NewEngine(
		appx.AddConfigurator(config),
	).RUN(context.Background())
}

var config = &Config{}

type Config struct{}

func (c *Config) Valid() bool {
	return false
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("xxxx.json"),
		configx.NewFileReader("xxxx.json"),
	}
}

func (c *Config) Execute() error {
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
  name: quanx-test
  host: localhost
  port:
    http: 8080
    grpc: 8081
```

#### nacos配置

nacos配置文件路径：conf/nacos.yaml，不使用nacos可不添加。

```yaml
address: "127.0.0.1:8848"     # string nacos服务地址,多个以英文逗号分割
username: "nacos"             # string 用户名
password: "nacos"             # string 密码
namespace: "demo"             # string 命名空间
group: "DEFAULT_GROUP"        # string 配置分组
mode: 2                       # int 模式（0-仅配置中心；1-仅服务发现；2-配置中心和服务发现）
```

#### 数据库配置

数据库配置文件路径：conf/database.yaml

```yaml
source: "default"             # string 数据源名称
client: "gorm"                # string 客户端选型
enable: false                 # bool 是否启用
type: "mysql"                 # string 数据库类型(mysql/postgres)
host: "127.0.0.1"             # string host
port: 5432                    # int 端口
username: "root"              # string 用户名
password: "root"              # string 密码
database: ""                  # string 数据库名
schema: ""                    # string 模式名（postgres）
maxIdleConns: 10              # int 最大空闲连接数
maxOpenConns: 10              # int 最大打开连接数
connMaxLifetime: 10           # int 连接存活时间(秒)
logLevel: debug               # string 日志级别
slowThreshold: 200            # int 慢查询阈值（毫秒）
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

#### 缓存配置

redis配置文件路径：conf/cache.yaml

```yaml
source: "default"             # string 数据源名称
client: "local"               # string 客户端选型（redis/local）
enable: true                  # bool 是否启用
address: "localhost"          # string 地址
password: ""                  # string 密码
database: 0                   # int 数据库
mode: 0                       # int 模式（0-单机；1-集群），默认单机模式
```

##### 多redis源

如果需要连接多个redis数据源，更新conf/redis.yaml配置文件内容修改为多配置即可

```yaml
- name: default
  client: local
  enable:
  address: 
  password: 
  database: 
  mode: 0
- name: redis_db1
  client: redis
  enable:
  address: 
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
package main

import (
	"context"
	"fmt"

	"github.com/go-xuan/quanx/appx"
	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

func main() {
	appx.NewEngine(
		appx.AddConfigurator(&demo{}),
	).RUN(context.Background())
}

type demo struct {
	Key1 int      `json:"key1" yaml:"key1"`
	Key2 string   `json:"key2" yaml:"key2"`
	Key3 []string `json:"key3" yaml:"key3"`
}

func (d *demo) Valid() bool {
	return d.Key1 > 0 && d.Key2 != "" && d.Key3 != nil
}

func (d *demo) Readers() []configx.Reader {
	return []configx.Reader{
		configx.NewFileReader("demo.yaml"),
		configx.NewFileReader("demo.json"),
		nacosx.NewReader("demo.yaml"),
	}
}

func (d *demo) Execute() error {
	// todo 完成配置读取后需要进行的操作
	fmt.Println(d.Key1)
	fmt.Println(d.Key2)
	fmt.Println(d.Key3)
	return nil
}


```


