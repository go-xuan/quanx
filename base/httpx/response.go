package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/go-xuan/quanx/base/errorx"
)

// Response 表示 HTTP 请求的响应结构
type Response struct {
	trace   string         // 追踪标识，用于日志或调试
	status  int            // HTTP 状态码
	body    []byte         // 响应体
	cookies []*http.Cookie // 响应中的 cookies
}

// StatusOK 检查响应状态码是否为 200 OK
func (r *Response) StatusOK() bool {
	return r.status == http.StatusOK
}

// ValidStatus 检查响应状态码是否在指定状态码列表中
// 如果未提供状态码列表，则默认检查是否为 200 OK
func (r *Response) ValidStatus(status ...int) bool {
	if len(status) == 0 {
		return r.StatusOK()
	}
	for _, v := range status {
		if v == r.status {
			return true
		}
	}
	return false
}

// NotNull 检查响应体是否不为空
func (r *Response) NotNull() bool {
	return r.body != nil && len(r.body) > 0
}

// Unmarshal 将响应体解析到指定的结构体中
// 如果响应体为空或解析失败，将返回错误
func (r *Response) Unmarshal(v any) error {
	if r.body == nil || len(r.body) == 0 {
		return errorx.New("response body is empty, cannot unmarshal")
	}
	if err := json.Unmarshal(r.body, v); err != nil {
		return errorx.Wrap(err, "json unmarshal error")
	}
	return nil
}

// Cookies 返回响应中的 cookies
func (r *Response) Cookies() []*http.Cookie {
	return r.cookies
}

// Body 返回响应体的字节切片
func (r *Response) Body() []byte {
	return r.body
}

// JsonResult 使用 gjson 解析响应体并返回解析结果
// 如果响应体为空，将返回空的 gjson.Result
func (r *Response) JsonResult() gjson.Result {
	return gjson.ParseBytes(r.body)
}
