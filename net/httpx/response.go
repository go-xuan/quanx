package httpx

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	code    int
	body    []byte
	cookies []*http.Cookie
}

func (r *Response) StatusOK() bool {
	return r.code == http.StatusOK
}

func (r *Response) Unmarshal(v any) error {
	if err := json.Unmarshal(r.body, v); err != nil {
		return err
	}
	return nil
}

func (r *Response) Cookies() []*http.Cookie {
	return r.cookies
}

func (r *Response) Body() []byte {
	return r.body
}
