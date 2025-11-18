package modelx

import (
	"strconv"
	"strings"
)

// PageReq 分页参数
type PageReq struct {
	PageNo   int `json:"pageNo"`   // 分页页码
	PageSize int `json:"pageSize"` // 分页大小
}

// PageTotal 计算分页数量
func (r *PageReq) PageTotal(total int) int {
	if total != 0 && r.PageSize != 0 {
		if total%r.PageSize > 0 {
			return (total / r.PageSize) + 1
		} else {
			return total / r.PageSize
		}
	}
	return 0
}

// Offset 计算offset
func (r *PageReq) Offset() int {
	if r.PageNo != 0 && r.PageSize != 0 {
		return (r.PageNo - 1) * r.PageSize
	}
	return 0
}

// PgPageSql 获取PG分页SQL
func (r *PageReq) PgPageSql() string {
	if r != nil && r.PageNo != 0 {
		sb := strings.Builder{}
		sb.WriteString(` limit `)
		sb.WriteString(strconv.Itoa(r.PageSize))
		sb.WriteString(` offset `)
		sb.WriteString(strconv.Itoa(r.Offset()))
		return sb.String()
	}
	return ""
}

// MysqlPageSql 获取Mysql分页SQL
func (r *PageReq) MysqlPageSql() string {
	if r != nil && r.PageNo != 0 {
		sb := strings.Builder{}
		sb.WriteString(` limit `)
		sb.WriteString(strconv.Itoa(r.Offset() + 1))
		sb.WriteString(` , `)
		sb.WriteString(strconv.Itoa(r.PageSize))
		return sb.String()
	}
	return ""
}

// BuildResp 封装分页结果
func (r *PageReq) BuildResp(rows any, total int64) *PageResp {
	resp := &PageResp{
		Total: total,
		Rows:  rows,
	}
	if r != nil {
		resp.PageNo = r.PageNo
		resp.PageSize = r.PageSize
		resp.PageTotal = r.PageTotal(int(total))
	}
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
