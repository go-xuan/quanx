package httpx

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-xuan/quanx/os/errorx"
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

func (r *Request) do(strategy ClientStrategy) (resp *Response, err error) {
	if r.url == "" {
		err = errorx.New("url is empty")
		return
	}
	var httpRequest *http.Request
	if httpRequest, err = http.NewRequest(r.method, r.url, r.body); err != nil {
		return
	}
	if r.headers != nil && len(r.headers) > 0 {
		if _, ok := r.headers["Content-Type"]; !ok {
			r.headers["Content-Type"] = "application/json"
		}
		for key, val := range r.headers {
			httpRequest.Header.Set(key, val)
		}
	}
	var httpResponse *http.Response
	if httpResponse, err = GetClient(strategy).client.Do(httpRequest); err != nil {
		return
	}
	resp = &Response{
		code:    httpResponse.StatusCode,
		cookies: httpResponse.Cookies(),
	}
	defer httpResponse.Body.Close()
	var body []byte
	if body, err = io.ReadAll(httpResponse.Body); err != nil {
		return
	}
	resp.body = body
	if r.debug {
		log.Printf("[debug] url: %s\n", r.url)
		log.Printf("[debug] body: %s\n", string(body))
	}
	return
}

func (r *Request) Do() (*Response, error) {
	return r.do(&HttpClientStrategy{})
}

func (r *Request) DoProxy(proxyUrl string) (*Response, error) {
	return r.do(&ProxyClientStrategy{Proxy: proxyUrl})
}

func (r *Request) DoHttps(crt string) (*Response, error) {
	return r.do(&HttpsClientStrategy{Crt: crt})
}

func (r *Request) DoHttpsProxy(proxyUrl, crt string) (*Response, error) {
	return r.do(&HttpsProxyClientStrategy{
		Proxy: proxyUrl,
		Crt:   crt,
	})
}
