package modelx

import (
	"mime/multipart"

	"gorm.io/gorm"
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
	field string      // 字段名
	Start interface{} // 开始时间
	End   interface{} // 结束时间
}

func (b Between) Where(db *gorm.DB) *gorm.DB {
	if b.Start != nil {
		db = db.Where(b.field+" >= ?", b.Start)
	}
	if b.End != nil {
		db = db.Where(b.field+" <= ?", b.End)
	}
	return db
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

// Query 分页参数
type Query struct {
	Page
	Keyword string  `json:"keyword"` // 关键字
	OrderBy []Order `json:"orderBy"` // 排序参数
}

func (q *Query) DoLike(db *gorm.DB, fields ...string) *gorm.DB {
	if q.Keyword != "" && len(fields) > 0 {
		likeKeyword := "%" + q.Keyword + "%"
		for _, field := range fields {
			db = db.Where(field+" LIKE ?", likeKeyword)
		}
	}
	return db
}

func (q *Query) DoPage(db *gorm.DB) *gorm.DB {
	if q.Page.PageSize > 0 {
		db = db.Limit(q.Page.PageSize).Offset(q.Page.Offset())
	}
	return db
}

func (q *Query) SetOrder(db *gorm.DB, def string) *gorm.DB {
	if q.OrderBy != nil {
		for _, order := range q.OrderBy {
			db = db.Order(order.Value())
		}
	} else {
		db = db.Order(def)
	}
	return db
}
