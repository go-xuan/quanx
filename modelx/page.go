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

// IsValid 是否有效
func (p Page) IsValid() bool {
	return p.PageNo > 0 || p.PageSize > 0
}

// PageTotal 计算分页数量
func (p Page) PageTotal(total int64) int {
	if total > 0 && p.PageSize > 0 {
		if total%int64(p.PageSize) > 0 {
			return int((total / int64(p.PageSize)) + 1)
		}
		return int(total / int64(p.PageSize))
	}
	return 0
}

// Offset 获取offset
func (p Page) Offset() int {
	if p.PageNo > 0 && p.PageSize > 0 {
		return (p.PageNo - 1) * p.PageSize
	}
	return 0
}

// Limit 获取limit
func (p Page) Limit() int {
	return p.PageSize
}

// PgPageSql 获取PG分页SQL
func (p Page) PgPageSql() string {
	if p.PageNo > 0 {
		sb := strings.Builder{}
		sb.WriteString(` limit `)
		sb.WriteString(strconv.Itoa(p.PageSize))
		sb.WriteString(` offset `)
		sb.WriteString(strconv.Itoa(p.Offset()))
		return sb.String()
	}
	return ""
}

// MysqlPageSql 获取Mysql分页SQL
func (p Page) MysqlPageSql() string {
	if p.PageNo > 0 {
		sb := strings.Builder{}
		sb.WriteString(` limit `)
		sb.WriteString(strconv.Itoa(p.Offset() + 1))
		sb.WriteString(` , `)
		sb.WriteString(strconv.Itoa(p.PageSize))
		return sb.String()
	}
	return ""
}

// BuildResp 封装分页结果
func (p Page) BuildResp(rows any, total int64) *PageResp {
	resp := &PageResp{
		Total: total,
		Rows:  rows,
	}
	resp.PageNo = p.PageNo
	resp.PageSize = p.PageSize
	resp.PageTotal = p.PageTotal(total)
	return resp
}

// PageResp 分页结果
type PageResp struct {
	PageNo    int   `json:"pageNo" comment:"当前页码"`
	PageSize  int   `json:"pageSize" comment:"当前页容量"`
	PageTotal int   `json:"pageTotal" comment:"总页数"`
	Total     int64 `json:"total" comment:"结果总数量"`
	Rows      any   `json:"rows" comment:"返回结果集"`
}
