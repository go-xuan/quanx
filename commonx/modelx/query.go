package modelx

import "strings"

type Id[T any] struct {
	Id T `form:"id" json:"id" binding:"required"`
}
type Ids[T any] struct {
	Ids []T `form:"ids" json:"ids" binding:"required"`
}

// 分页参数
type Query struct {
	Keyword string `json:"keyword"` // 关键字
	Order   Orders `json:"order"`   // 排序参数
	TimeRange
	Page
}

type Orders []*Order
type Order struct {
	Column string `json:"column"` // 排序字段
	Type   string `json:"type"`   // 排序方式(asc/desc)
}

// 计算分页数量
func (params Orders) GetOrderBySql() (orderBy string) {
	if params != nil && len(params) > 0 {
		s := strings.Builder{}
		s.WriteString(` order by `)
		for i, param := range params {
			if i > 1 {
				s.WriteString(",")
			}
			s.WriteString(param.Column)
			s.WriteString(" ")
			s.WriteString(param.Type)
		}
		orderBy = s.String()
	}
	return
}

// 时间范围
type TimeRange struct {
	StartTime string `json:"startTime"` // 开始时间
	EndTime   string `json:"endTime"`   // 结束时间
}
