package respx

import (
	"github.com/go-xuan/quanx/app/modelx"
)

// PageResponse 分页结果
type PageResponse struct {
	PageNo    int64 `json:"pageNo" comment:"当前页码"`
	PageSize  int64 `json:"pageSize" comment:"当前页容量"`
	PageTotal int64 `json:"pageTotal" comment:"总页数"`
	Total     int64 `json:"total" comment:"结果总数量"`
	Rows      any   `json:"rows" comment:"返回结果集"`
}

// BuildPageResp 封装分页结果
func BuildPageResp(page *modelx.Page, rows any, total int64) *PageResponse {
	return &PageResponse{
		PageNo:    page.PageNo,
		PageSize:  page.PageSize,
		PageTotal: page.PageTotal(total),
		Total:     total,
		Rows:      rows,
	}
}
