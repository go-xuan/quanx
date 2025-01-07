package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/go-xuan/quanx/os/errorx"
)

type Response struct {
	status  int
	body    []byte
	cookies []*http.Cookie
}

func (r *Response) StatusOK() bool {
	return r.status == http.StatusOK
}

func (r *Response) ValidStatus(status ...int) bool {
	if len(status) == 0 {
		return r.StatusOK()
	} else {
		for _, v := range status {
			if v == r.status {
				return true
			}
		}
	}
	return false
}

func (r *Response) NotNull() bool {
	return r.body != nil && len(r.body) > 0
}

func (r *Response) Unmarshal(v any) error {
	if err := json.Unmarshal(r.body, v); err != nil {
		return errorx.Wrap(err, "json unmarshal error")
	}
	return nil
}

func (r *Response) Cookies() []*http.Cookie {
	return r.cookies
}

func (r *Response) Body() []byte {
	return r.body
}

func (r *Response) JsonResult() gjson.Result {
	return gjson.ParseBytes(r.body)
}
