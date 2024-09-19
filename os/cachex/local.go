package cachex

import (
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/go-xuan/quanx/app/configx"
)

var localCache *cache.Cache

func NewLocalCache() *cache.Cache {
	if localCache == nil {
		localCache = cache.New(time.Duration(-1), time.Duration(-1))
	}
	return localCache
}

type LocalCache struct {
	Expiration int    `json:"expiration" yaml:"expiration" default:"-1"` // 默认缓存有效期，-1表示无期
	Interval   int    `json:"interval" yaml:"interval" default:"-1"`     // 缓存清理间隔，-1表示永不清理
	File       string `json:"file" yaml:"file"`                          // 缓存加载文件
}

func (*LocalCache) ID() string {
	return "local-cache"
}

func (*LocalCache) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "local_cache.yaml",
		NacosDataId: "local_cache.yaml",
		Listen:      false,
	}
}

func (conf *LocalCache) Execute() error {
	localCache = cache.New(time.Duration(conf.Expiration), time.Duration(conf.Interval))
	if conf.File != "" {
		return localCache.LoadFile(conf.File)
	}
	return nil
}
