package httpx

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	POST          = "POST"
	GET           = "GET"
	PUT           = "PUT"
	DELETE        = "DELETE"
	HttpPrefix    = "http://"
	ContentType   = "Content-Type"
	Authorization = "Authorization"
	JsonContent   = "application/json"
	FormContent   = "application/x-www-form-urlencoded"
	GBKContent    = "application/javascript;charset=GBK"
)

func RequestHttp(method string, host string, header map[string]string, param interface{}) (result string, err error) {
	var paramJson []byte
	paramJson, err = json.Marshal(param)
	if err != nil {
		log.Error(err)
		return
	}
	var req *http.Request
	req, err = http.NewRequest(method, host, bytes.NewBuffer(paramJson))
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set(ContentType, JsonContent)
	for key, val := range header {
		req.Header.Set(key, val)
	}
	var resp *http.Response
	resp, err = getHttpClient().Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer func(Body *http.Response) {
		_ = resp.Body.Close()
	}(resp)
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	if len(resp.Header[ContentType]) > 0 {
		if resp.Header[ContentType][0] == GBKContent {
			result = DecodeBody(string(body))
		} else {
			result = string(body)
		}
	} else {
		result = string(body)
	}
	return
}

// 转码
func DecodeBody(src string) (dst string) {
	data, err := io.ReadAll(transform.NewReader(bytes.NewBuffer([]byte(src)), simplifiedchinese.GBK.NewDecoder()))
	if err == nil {
		dst = string(data)
	}
	return
}

// map转为Url
func MapToUrl(params map[string]interface{}) (s string) {
	isFirst := false
	for k, v := range params {
		if isFirst {
			s = s + "&"
		}
		s = s + k + "=" + url.QueryEscape(typeSwitcher(v))
		isFirst = true
	}
	return
}

func typeSwitcher(t interface{}) string {
	switch v := t.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	case int64:
		return strconv.Itoa(int(v))
	case []string:
		return "typeArray"
	case map[string]interface{}:
		return "typeMap"
	default:
		return ""
	}
}
