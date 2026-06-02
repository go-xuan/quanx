package modelx

import (
	"mime/multipart"
)

type Id[T int | uint | int64 | uint64 | string] struct {
	Id T `form:"id" json:"id" binding:"required"`
}

type Ids[T int | uint | int64 | uint64 | string] struct {
	Ids []T `form:"ids" json:"ids" binding:"required"`
}

type File struct {
	File *multipart.FileHeader `form:"file"`
}

// Between 条件范围
type Between struct {
	Field string // 字段名
	Start any    // 开始值
	End   any    // 结束值
}

// Order 排序参数
type Order struct {
	Column string `json:"column"` // 排序字段
	Desc   bool   `json:"desc"`   // 是否降序排序
}

func (o Order) Value() string {
	if o.Desc {
		return o.Column + " desc"
	}
	return o.Column + " asc"
}

// Query 分页查询参数
type Query struct {
	Page
	Keyword string  `json:"keyword"` // 关键字
	OrderBy []Order `json:"orderBy"` // 排序参数
}
