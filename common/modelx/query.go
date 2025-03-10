package modelx

import (
	"mime/multipart"
	"strings"
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

// Query 分页参数
type Query struct {
	Keyword string `json:"keyword"` // 关键字
	Order   Orders `json:"order"`   // 排序参数
	TimeRange
	Page
}

// TimeRange 时间范围
type TimeRange struct {
	StartTime string `json:"startTime"` // 开始时间
	EndTime   string `json:"endTime"`   // 结束时间
}

// Orders 排序
type Orders []*Order
type Order struct {
	Column string `json:"column"` // 排序字段
	Type   string `json:"type"`   // 排序方式(asc/desc)
}

// GetOrderBySql 计算分页数量
func (orders Orders) GetOrderBySql() string {
	if orders != nil && len(orders) > 0 {
		s := strings.Builder{}
		s.WriteString(` order by `)
		for i, order := range orders {
			if i > 1 {
				s.WriteString(",")
			}
			s.WriteString(order.Column)
			s.WriteString(" ")
			s.WriteString(order.Type)
		}
		return s.String()
	}
	return ""
}
