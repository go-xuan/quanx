package modelx

import (
	"strconv"
	"strings"
)

// 分页参数
type Page struct {
	PageNo   int64 `json:"pageNo" comment:"分页页码"`
	PageSize int64 `json:"pageSize" comment:"分页大小"`
}

// 计算分页数量
func (page *Page) PageTotal(total int64) int64 {
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

// 计算offset
func (page *Page) Offset() int64 {
	if page.PageNo != 0 && page.PageSize != 0 {
		return (page.PageNo - 1) * page.PageSize
	} else {
		return 0
	}
}

// 获取PG分页SQL
func (page *Page) PgPageSql() string {
	pageSql := strings.Builder{}
	if page != nil && page.PageNo != 0 {
		pageSql.WriteString(` limit `)
		pageSql.WriteString(strconv.FormatInt(page.PageSize, 10))
		pageSql.WriteString(` offset `)
		pageSql.WriteString(strconv.FormatInt(page.Offset(), 10))
	}
	return pageSql.String()
}

// 获取Mysql分页SQL
func (page *Page) MysqlPageSql() string {
	pageSql := strings.Builder{}
	if page != nil && page.PageNo != 0 {
		pageSql.WriteString(` limit `)
		pageSql.WriteString(strconv.FormatInt(page.Offset()+1, 10))
		pageSql.WriteString(` , `)
		pageSql.WriteString(strconv.FormatInt(page.PageSize, 10))
	}
	return pageSql.String()
}
