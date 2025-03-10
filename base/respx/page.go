package respx

import (
	"github.com/go-xuan/quanx/common/modelx"
)

// PageResponse 分页结果
type PageResponse struct {
	PageNo    int   `json:"pageNo" comment:"当前页码"`
	PageSize  int   `json:"pageSize" comment:"当前页容量"`
	PageTotal int   `json:"pageTotal" comment:"总页数"`
	Total     int64 `json:"total" comment:"结果总数量"`
	Rows      any   `json:"rows" comment:"返回结果集"`
}

// BuildPageResp 封装分页结果
func BuildPageResp(page *modelx.Page, rows any, total int64) *PageResponse {
	resp := &PageResponse{
		Total: total,
		Rows:  rows,
	}
	if page != nil {
		resp.PageNo = page.PageNo
		resp.PageSize = page.PageSize
		resp.PageTotal = page.PageTotal(int(total))
	}
	return resp
}
