package dbx

import (
	"gorm.io/gorm"

	"github.com/go-xuan/quanx/modelx"
)

// BetweenWhere 将 Between 条件应用到 gorm 查询
func BetweenWhere(db *gorm.DB, b modelx.Between) *gorm.DB {
	if b.Start != nil {
		db = db.Where(b.Field+" >= ?", b.Start)
	}
	if b.End != nil {
		db = db.Where(b.Field+" <= ?", b.End)
	}
	return db
}

// QueryLike 将 Query 的关键词搜索应用到 gorm 查询
func QueryLike(db *gorm.DB, q *modelx.Query, fields ...string) *gorm.DB {
	if q.Keyword != "" && len(fields) > 0 {
		likeKeyword := "%" + q.Keyword + "%"
		for _, field := range fields {
			db = db.Where(field+" LIKE ?", likeKeyword)
		}
	}
	return db
}

// QueryPage 将 Query 的分页参数应用到 gorm 查询
func QueryPage(db *gorm.DB, q *modelx.Query) *gorm.DB {
	if q.PageSize > 0 {
		db = db.Limit(q.PageSize).Offset(q.Offset())
	}
	return db
}

// QueryOrder 将 Query 的排序参数应用到 gorm 查询
func QueryOrder(db *gorm.DB, q *modelx.Query, def string) *gorm.DB {
	if q.OrderBy != nil {
		for _, order := range q.OrderBy {
			db = db.Order(order.Value())
		}
	} else {
		db = db.Order(def)
	}
	return db
}
