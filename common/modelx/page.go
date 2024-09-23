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
	pageSql := strings.Builder{}
	if page != nil && page.PageNo != 0 {
		pageSql.WriteString(` limit `)
		pageSql.WriteString(strconv.Itoa(page.PageSize))
		pageSql.WriteString(` offset `)
		pageSql.WriteString(strconv.Itoa(page.Offset()))
	}
	return pageSql.String()
}

// MysqlPageSql 获取Mysql分页SQL
func (page *Page) MysqlPageSql() string {
	pageSql := strings.Builder{}
	if page != nil && page.PageNo != 0 {
		pageSql.WriteString(` limit `)
		pageSql.WriteString(strconv.Itoa(page.Offset() + 1))
		pageSql.WriteString(` , `)
		pageSql.WriteString(strconv.Itoa(page.PageSize))
	}
	return pageSql.String()
}
