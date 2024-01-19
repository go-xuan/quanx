package httpx

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Request struct {
	url    string
	method string
	header map[string]string
	body   interface{}
	form   url.Values
}

func New(method ...string) *Request {
	if len(method) > 0 {
		return &Request{method: method[0]}
	}
	return Get()
}

func Get() *Request {
	return &Request{method: GET}
}

func Post() *Request {
	return &Request{method: POST}
}

func (r *Request) Url(url string) *Request {
	r.url = url
	return r
}

func (r *Request) Method(method string) *Request {
	r.method = method
	return r
}
func (r *Request) Header(header map[string]string) *Request {
	r.header = header
	return r
}
func (r *Request) Body(body interface{}) *Request {
	r.body = body
	return r
}

func (r *Request) Form(form url.Values) *Request {
	r.form = form
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	if r.header != nil && len(r.header) > 0 {
		r.header[key] = value
	} else {
		var header = make(map[string]string)
		header[key] = value
		r.header = header
	}
	return r
}

func (r *Request) Authorization(token string) *Request {
	r.SetHeader("Authorization", token)
	return r
}

func (r *Request) Cookie(cookie string) *Request {
	r.SetHeader("Cookie", cookie)
	return r
}

func (r *Request) Do(modeAndParam ...string) (res []byte, err error) {
	var body io.Reader
	var contentType string
	if r.form != nil {
		r.method = POST
		contentType = "application/x-www-form-urlencoded"
		body = strings.NewReader(r.form.Encode())
	} else {
		contentType = "application/json"
		marshal, _ := json.Marshal(r.body)
		body = bytes.NewReader(marshal)
	}
	var req *http.Request
	req, err = http.NewRequest(r.url, r.url, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", contentType)
	if r.header != nil && len(r.header) > 0 {
		for key, val := range r.header {
			req.Header.Set(key, val)
		}
	}
	// 切换http客户端
	var resp *http.Response
	resp, err = SwitchClient(modeAndParam...).Do(req)
	if err != nil {
		return
	}
	defer func(resp *http.Response) {
		_ = resp.Body.Close()
	}(resp)
	res, err = io.ReadAll(resp.Body)
	return
}

// map转为Url
func MapToUrl(params map[string]interface{}) (s string) {
	sb := strings.Builder{}
	for k, v := range params {
		sb.WriteString("&")
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(typeSwitcher(v))
	}
	return sb.String()[1:]
}

func typeSwitcher(t interface{}) string {
	switch v := t.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}
