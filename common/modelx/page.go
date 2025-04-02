package modelx

import (
	"strconv"
	"strings"
)

// Page 分页参数
type Page struct {
	PageNo   int `json:"pageNo"`   // 分页页码
	PageSize int `json:"pageSize"` // 分页大小
}

type Order struct {
	Column string `json:"column"` // 排序字段
	Type   string `json:"type"`   // 排序方式(asc/desc)
}

// PageTotal 计算分页数量
func (page *Page) PageTotal(total int) int {
	if total != 0 && page.PageSize != 0 {
		if total%page.PageSize > 0 {
			return (total / page.PageSize) + 1
		} else {
			return total / page.PageSize
		}
	} else {
		return 0
	}
}

// Offset 计算offset
func (page *Page) Offset() int {
	if page.PageNo != 0 && page.PageSize != 0 {
		return (page.PageNo - 1) * page.PageSize
	} else {
		return 0
	}
}

// PgPageSql 获取PG分页SQL
func (page *Page) PgPageSql() string {
	if page != nil && page.PageNo != 0 {
		sb := strings.Builder{}
		sb.WriteString(` limit `)
		sb.WriteString(strconv.Itoa(page.PageSize))
		sb.WriteString(` offset `)
		sb.WriteString(strconv.Itoa(page.Offset()))
		return sb.String()
	}
	return ""
}

// MysqlPageSql 获取Mysql分页SQL
func (page *Page) MysqlPageSql() string {
	if page != nil && page.PageNo != 0 {
		sb := strings.Builder{}
		sb.WriteString(` limit `)
		sb.WriteString(strconv.Itoa(page.Offset() + 1))
		sb.WriteString(` , `)
		sb.WriteString(strconv.Itoa(page.PageSize))
		return sb.String()
	}
	return ""
}
