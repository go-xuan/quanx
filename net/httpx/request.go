package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Request struct {
	method  string
	url     string
	headers map[string]string
	body    io.Reader
	debug   bool
}

func Method(method string, url string) *Request {
	return &Request{
		method: method,
		url:    url,
	}
}

func Get(url string) *Request {
	return &Request{
		method: "GET",
		url:    url,
	}
}

func Post(url string) *Request {
	return &Request{
		method: "POST",
		url:    url,
	}
}

func (r *Request) Params(params map[string]string) *Request {
	sb := strings.Builder{}
	for k, v := range params {
		sb.WriteString("&")
		sb.WriteString(url.QueryEscape(k))
		sb.WriteString("=")
		sb.WriteString(url.QueryEscape(v))
	}
	r.url = r.url + "?" + sb.String()[1:]
	return r
}

func (r *Request) Debug() *Request {
	r.debug = true
	return r
}

func (r *Request) Body(body any) *Request {
	marshal, _ := json.Marshal(body)
	r.body = bytes.NewReader(marshal)
	r.SetHeader("Content-Type", "application/json")
	return r
}

func (r *Request) Form(form url.Values) *Request {
	r.body = strings.NewReader(form.Encode())
	r.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	return r
}

func (r *Request) Headers(headers map[string]string) *Request {
	r.headers = headers
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	if r.headers != nil && len(r.headers) > 0 {
		r.headers[key] = value
	} else {
		var header = make(map[string]string)
		header[key] = value
		r.headers = header
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

func (r *Request) do(strategy ClientStrategy) (res *Response, err error) {
	if r.url == "" {
		err = errors.New("url is empty")
		return
	}

	var req *http.Request
	var resp *http.Response
	if req, err = http.NewRequest(r.method, r.url, r.body); err != nil {
		return
	}
	if r.headers != nil && len(r.headers) > 0 {
		if _, ok := r.headers["Content-Type"]; !ok {
			r.headers["Content-Type"] = "application/json"
		}
		for key, val := range r.headers {
			req.Header.Set(key, val)
		}
	}
	if resp, err = GetClient(strategy).client.Do(req); err != nil {
		return
	}
	res = &Response{
		code:    resp.StatusCode,
		cookies: resp.Cookies(),
	}
	defer resp.Body.Close()
	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return
	}
	res.body = body
	if r.debug {
		log.Printf("[debug] url: %s\n", r.url)
		log.Printf("[debug] body: %s\n", string(body))
	}
	return
}

func (r *Request) Do() (response *Response, err error) {
	return r.do(&HttpClientStrategy{})
}

func (r *Request) DoProxy(proxyUrl string) (resp *Response, err error) {
	return r.do(&ProxyClientStrategy{
		Proxy: proxyUrl,
	})
}

func (r *Request) DoHttps(crt string) (resp *Response, err error) {
	return r.do(&HttpsClientStrategy{
		Crt: crt,
	})
}

func (r *Request) DoHttpsProxy(proxyUrl, crt string) (resp *Response, err error) {
	return r.do(&HttpsProxyClientStrategy{
		Proxy: proxyUrl,
		Crt:   crt,
	})
}
