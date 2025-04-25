package modelx

import (
	"gorm.io/gorm"
	"mime/multipart"
)

type Id[T int | int64 | string] struct {
	Id T `form:"id" json:"id" binding:"required"`
}

type Ids[T int | int64 | string] struct {
	Ids []T `form:"ids" json:"ids" binding:"required"`
}

type File struct {
	File *multipart.FileHeader `form:"file"`
}

// Between 条件范围
type Between struct {
	Start interface{} `json:"start"` // 开始时间
	End   interface{} `json:"end"`   // 结束时间
}

func (b Between) DoBetween(db *gorm.DB, field string) *gorm.DB {
	if b.Start != nil {
		db = db.Where(field+" >= ?", b.Start)
	}
	if b.End != nil {
		db = db.Where(field+" <= ?", b.End)
	}
	return db
}

// Query 分页参数
type Query struct {
	Keyword string   `json:"keyword"` // 关键字
	OrderBy []*Order `json:"orderBy"` // 排序参数
	*Page
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
	if q.Page != nil && q.Page.PageSize > 0 {
		db = db.Limit(q.Page.PageSize).Offset(q.Page.Offset())
	}
	return db
}

func (q *Query) DoOrder(db *gorm.DB, def string) *gorm.DB {
	if q.OrderBy != nil {
		for _, order := range q.OrderBy {
			db = db.Order(order.Column + " " + order.Type)
		}
	} else {
		db = db.Order(def)
	}
	return db
}
