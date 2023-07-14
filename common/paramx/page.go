package paramx

import (
	"strconv"
	"strings"
)

// 分页参数
type Page struct {
	CurrPage int64 `json:"currPage" comment:"当前页码"`
	PageSize int64 `json:"pageSize" comment:"当前页数量"`
}

// 计算分页数量
func (page *Page) GetPageTotal(total int64) int64 {
	if total != 0 && page.PageSize != 0 {
		if total%page.PageSize > 0 {
			return total/page.PageSize + 1
		} else {
			return total / page.PageSize
		}
	} else {
		return 0
	}
}

// 计算offset(pg/oracle)
func (page *Page) GetOffset() int64 {
	if page.CurrPage != 0 && page.PageSize != 0 {
		return (page.CurrPage - 1) * page.PageSize
	} else {
		return 0
	}
}

// 获取PG分页SQL
func (page *Page) GetPgPageSql() string {
	pageSql := strings.Builder{}
	if page != nil && page.CurrPage != 0 {
		pageSql.WriteString(` limit `)
		pageSql.WriteString(strconv.FormatInt(page.PageSize, 10))
		pageSql.WriteString(` offset `)
		pageSql.WriteString(strconv.FormatInt(page.GetOffset(), 10))
	}
	return pageSql.String()
}

// 获取Mysql分页SQL
func (page *Page) GetMysqlPageSql() string {
	pageSql := strings.Builder{}
	if page != nil && page.CurrPage != 0 {
		pageSql.WriteString(` limit `)
		pageSql.WriteString(strconv.FormatInt(page.GetOffset()+1, 10))
		pageSql.WriteString(` , `)
		pageSql.WriteString(strconv.FormatInt(page.PageSize, 10))
	}
	return pageSql.String()
}
