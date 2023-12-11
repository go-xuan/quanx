package quanx

import (
	"github.com/go-xuan/quanx/gormx"
	"github.com/go-xuan/quanx/logx"
	"github.com/go-xuan/quanx/nacosx"
	"github.com/go-xuan/quanx/redisx"
)

type X[V any] interface {
	Format() string
	Init()
}

var (
	LogX   X[[]*logx.Config]
	NacosX X[*nacosx.Config]
	GormX  X[[]*gormx.Config]
	RedisX X[[]*redisx.Config]
	//ElasticX   X[*elasticx.Config]
	//MinioX     X[*miniox.Config]
	//MongoX     X[*mongox.Config]
	//HugegraphX X[*hugegraphx.Config]
)
