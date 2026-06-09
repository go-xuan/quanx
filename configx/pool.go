package configx

import (
	"errors"

	"github.com/go-xuan/typex"
	"github.com/go-xuan/utilx/errorx"
)

// Pool 通用客户端池，提供多数据源客户端管理能力
type Pool[C any] struct {
	enum *typex.Enum[string, C]
}

// NewPool 创建客户端池
func NewPool[C any]() *Pool[C] {
	return &Pool[C]{}
}

// Initialized 是否已初始化
func (p *Pool[C]) Initialized() bool {
	return p.enum != nil && p.enum.Len() > 0
}

// Add 添加客户端，首个添加的客户端同时设为 default
func (p *Pool[C]) Add(source string, client C) {
	if p.enum == nil {
		p.enum = typex.NewStringEnum[C]()
		p.enum.Add("default", client)
	}
	p.enum.Add(source, client)
}

// Get 获取客户端，未指定 source 时返回 default
func (p *Pool[C]) Get(source ...string) C {
	if len(source) > 0 && source[0] != "" {
		if client := p.enum.Get(source[0]); any(client) != nil {
			return client
		}
	}
	return p.enum.Get("default")
}

// Range 遍历所有客户端
func (p *Pool[C]) Range(f func(source string, client C) bool) {
	if p.enum != nil {
		p.enum.Range(f)
	}
}

// Len 返回客户端数量
func (p *Pool[C]) Len() int {
	if p.enum == nil {
		return 0
	}
	return p.enum.Len()
}

// Close 遍历关闭所有客户端
func (p *Pool[C]) Close(closeFn func(C) error) error {
	var errs []error
	p.Range(func(source string, client C) bool {
		if e := closeFn(client); e != nil {
			errs = append(errs, errorx.Wrap(e, "close client failed"))
		}
		return true
	})
	return errors.Join(errs...)
}
