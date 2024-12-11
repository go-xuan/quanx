package httpx

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/os/errorx"
)

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

type Request struct {
	method  string
	url     string
	headers map[string]string
	body    io.Reader
	debug   bool
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

func (r *Request) SetAuthorization(token string) *Request {
	r.SetHeader("Authorization", token)
	return r
}

func (r *Request) SetCookie(cookie string) *Request {
	r.SetHeader("Cookie", cookie)
	return r
}

func (r *Request) Do(strategy ...ClientStrategy) (*Response, error) {
	if r.url == "" {
		return nil, errorx.New("url is empty")
	}
	httpRequest, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		return nil, errorx.Wrap(err, "http.NewRequest error")
	}
	if r.headers != nil && len(r.headers) > 0 {
		if _, ok := r.headers["Content-AndOr"]; !ok {
			r.headers["Content-AndOr"] = "application/json"
		}
		for key, val := range r.headers {
			httpRequest.Header.Set(key, val)
		}
	}
	var httpResponse *http.Response
	if httpResponse, err = GetClient(strategy...).HttpClient().Do(httpRequest); err != nil {
		return nil, errorx.Wrap(err, "client.Do error")
	}
	resp := &Response{
		code:    httpResponse.StatusCode,
		cookies: httpResponse.Cookies(),
	}
	defer httpResponse.Body.Close()
	var body []byte
	if body, err = io.ReadAll(httpResponse.Body); err != nil {
		return resp, errorx.Wrap(err, "http.Response.Body read error")
	}
	resp.body = body
	if r.debug {
		log.Printf("[debug] url: %s", r.url)
		log.Printf("[debug] body: %s", string(body))
	}
	return resp, nil
}

func (r *Request) DoProxy(proxyUrl string) (*Response, error) {
	return r.Do(&ProxyClientStrategy{Proxy: proxyUrl})
}

func (r *Request) DoHttps(crt string) (*Response, error) {
	return r.Do(&HttpsClientStrategy{Crt: crt})
}

func (r *Request) DoHttpsProxy(proxyUrl, crt string) (*Response, error) {
	return r.Do(&HttpsProxyClientStrategy{
		Proxy: proxyUrl,
		Crt:   crt,
	})
}
