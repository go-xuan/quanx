package httpx

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/errorx"
)

func Method(method string, url string) *Request {
	return &Request{
		method: method,
		url:    url,
	}
}

func Get(url string) *Request {
	return &Request{
		method: http.MethodGet,
		url:    url,
	}
}

func Post(url string) *Request {
	return &Request{
		method: http.MethodPost,
		url:    url,
	}
}

type Request struct {
	trace   string
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

func (r *Request) Trace(trace string) *Request {
	r.trace = trace
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

func (r *Request) Do(uses ...ClientUse) (*Response, error) {
	if r.url == "" {
		return &Response{trace: r.trace}, errorx.New("url is empty")
	}
	request, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		return &Response{trace: r.trace}, errorx.Wrap(err, "new request error")
	}
	if r.headers != nil && len(r.headers) > 0 {
		if _, ok := r.headers["Content-Type"]; !ok {
			r.headers["Content-Type"] = "application/json"
		}
		for key, val := range r.headers {
			request.Header.Set(key, val)
		}
	}
	var response *Response
	if response, err = GetClient(uses...).Do(request); err != nil {
		return &Response{trace: r.trace}, errorx.Wrap(err, "do request error")
	}
	response.trace = r.trace
	if r.debug {
		logger := log.WithField("httpx_debug", true)
		if trace := r.trace; trace != "" {
			logger = logger.WithField("trace", trace)
		}
		logger.Printf("request_url: %s", r.url)
		logger.Printf("response_body: %s", string(response.Body()))
	}
	return response, nil
}
