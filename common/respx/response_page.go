package respx

import "github.com/go-xuan/quanx/common/modelx"

// 分页结果
type PageResponse struct {
	PageNo    int64       `json:"pageNo" comment:"当前页码"`
	PageSize  int64       `json:"pageSize" comment:"当前页容量"`
	PageTotal int64       `json:"pageTotal" comment:"总页数"`
	Total     int64       `json:"total" comment:"结果总数量"`
	Rows      interface{} `json:"rows" comment:"返回结果集"`
}

// 封装分页结果
func BuildPageResp(page *modelx.Page, rows interface{}, total int64) *PageResponse {
	var pageData = PageResponse{
		PageNo:    page.PageNo,
		PageSize:  page.PageSize,
		PageTotal: page.PageTotal(total),
		Total:     total,
		Rows:      rows,
	}
	return &pageData
}
