package respx

import "github.com/quanxiaoxuan/quanx/common/modelx"

// 分页结果
type PageResponse struct {
	PageNo    int64       `json:"pageNo" comment:"当前页码"`
	PageSize  int64       `json:"pageSize" comment:"当前页数量"`
	PageTotal int64       `json:"pageTotal" comment:"总页数"`
	Total     int64       `json:"total" comment:"结果总数量"`
	List      interface{} `json:"list" comment:"返回结果集"`
}

// 封装分页结果
func BuildPageResp(page *modelx.Page, data interface{}, total int64) *PageResponse {
	var pageData = PageResponse{
		PageNo:    page.PageNo,
		PageSize:  page.PageSize,
		PageTotal: page.Total(total),
		Total:     total,
		List:      data,
	}
	return &pageData
}
