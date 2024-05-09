package cachex

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/core/confx"
	"github.com/go-xuan/quanx/db/redisx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

// 缓存配置
type MultiCache []*Cache

type Cache struct {
	Source  string `json:"source" yaml:"source" default:"default"`   // 缓存存储数据源名称
	Prefix  string `json:"prefix" yaml:"prefix" default:"default"`   // 缓存KEY前缀前缀
	Marshal string `json:"marshal" yaml:"marshal" default:"msgpack"` // 序列化方案
}

func (c *Cache) Theme() string {
	return "Cache"
}

func (c *Cache) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "cache.yaml",
		NacosDataId: "cache.yaml",
		Listen:      false,
	}
}

func (c *Cache) Run() (err error) {
	if err = anyx.SetDefaultValue(c); err != nil {
		return
	}
	var client = c.CacheClient()
	handler = &Handler{
		Multi:     false,
		Client:    client,
		ClientMap: make(map[string]*CacheClient),
	}
	handler.ClientMap[c.Source] = client
	log.Info("Cache Init Successful : ", c.ToString())
	return
}

func Default() *Cache {
	return &Cache{
		Source:  constx.DefaultSourceName,
		Prefix:  "cache",
		Marshal: marshalx.Json,
	}
}

// 配置信息格式化
func (MultiCache) Theme() string {
	return "Cache"
}

// 配置文件读取
func (MultiCache) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "cache.yaml",
		NacosDataId: "cache.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (m MultiCache) Run() error {
	handler = &Handler{
		Multi:     true,
		ClientMap: make(map[string]*CacheClient),
	}
	multi := anyx.IfZero(m, MultiCache{Default()})
	for i, c := range multi {
		var client = c.CacheClient()
		handler.ClientMap[c.Source] = client
		if i == 0 || c.Source == constx.DefaultSourceName {
			handler.Client = client
		}
		log.Info("Cache Init Successful : ", c.ToString())
	}
	return nil
}

// 配置信息格式化
func (c *Cache) ToString() string {
	return fmt.Sprintf("name=%s prefix=%s marshal=%s", c.Source, c.Prefix, c.Marshal)
}

func (c *Cache) CacheClient() *CacheClient {
	return &CacheClient{
		cache:     c,
		client:    redisx.Client(c.Source),
		unmarshal: marshalx.NewCase(c.Marshal).Unmarshal,
	}
}

func (c *Cache) AddPrefix(keys []string) []string {
	var newKeys []string
	if len(keys) > 0 {
		for _, key := range keys {
			newKeys = append(newKeys, c.Prefix+key)
		}
	}
	return newKeys
}
