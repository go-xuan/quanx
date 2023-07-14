package paramx

import "strings"

// 排序
type OrderList []*Order
type Order struct {
	Field string `json:"field" comment:"排序字段"`
	Type  string `json:"type" comment:"排序方式(asc/desc)"`
}

// 计算分页数量
func (params OrderList) GetOrderBySql() (orderBy string) {
	if params != nil && len(params) > 0 {
		s := strings.Builder{}
		s.WriteString(` order by `)
		for i, param := range params {
			if i > 1 {
				s.WriteString(",")
			}
			s.WriteString(param.Field)
			s.WriteString(" ")
			s.WriteString(param.Type)
		}
		orderBy = s.String()
	}
	return
}
