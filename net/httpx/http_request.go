package httpx

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-xuan/quanx/types/anyx"
)

type Request struct {
	method string
	url    string
	header map[string]string
	body   any
	form   url.Values
	client *Client
}

func (r *Request) Https(crt string) *Request {
	r.client = newHttpsClient(crt)
	return r
}

func Method(method ...string) *Request {
	return &Request{method: anyx.Default(method, GET)}
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

func (r *Request) Body(body any) *Request {
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

func (r *Request) DoHttp() (res []byte, err error) {
	return r.Do()
}

func (r *Request) DoProxy(proxyUrl string) (res []byte, err error) {
	return r.Do(Proxy, proxyUrl)
}

func (r *Request) DoHttps(crt string) (res []byte, err error) {
	return r.Do(Https, crt)
}

func (r *Request) DoHttpsProxy(proxyUrl, crt string) (res []byte, err error) {
	return r.Do(HttpsProxy, proxyUrl, crt)
}

func (r *Request) Do(modeAndParam ...string) (res []byte, err error) {
	var body io.Reader
	if r.form != nil {
		r.method = POST
		r.SetHeader("Content-Type", "application/x-www-form-urlencoded")
		body = strings.NewReader(r.form.Encode())
	} else if r.body != nil {
		r.SetHeader("Content-Type", "application/json")
		marshal, _ := json.Marshal(r.body)
		body = bytes.NewReader(marshal)
	}
	var req *http.Request
	if req, err = http.NewRequest(r.method, r.url, body); err != nil {
		return
	}
	if r.header != nil && len(r.header) > 0 {
		for key, val := range r.header {
			req.Header.Set(key, val)
		}
	}
	// 切换http客户端
	var resp *http.Response
	if resp, err = SwitchClient(modeAndParam...).Do(req); err != nil {
		return
	}
	var reader = resp.Body
	defer reader.Close()
	return io.ReadAll(reader)
}

// map转为Url
func MapToUrl(params map[string]any) (s string) {
	sb := strings.Builder{}
	for k, v := range params {
		sb.WriteString("&")
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(typeSwitcher(v))
	}
	return sb.String()[1:]
}

func typeSwitcher(v any) string {
	switch t := v.(type) {
	case int:
		return strconv.Itoa(t)
	case string:
		return t
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return ""
	}
}
