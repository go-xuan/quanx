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
	Post            = "POST"
	Get             = "GET"
	Put             = "PUT"
	Delete          = "DELETE"
	HttpPrefix      = "http://"
	HttpsPrefix     = "https://"
	ContentTypeKey  = "Content-Type"
	ContentTypeJson = "application/json"
	ContentTypeGBK  = "application/javascript;charset=GBK"
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
	req.Header.Set(ContentTypeKey, ContentTypeJson)
	for key, val := range header {
		req.Header.Set(key, val)
	}
	var resp *http.Response
	resp, err = getHttpClient().Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	if len(resp.Header[ContentTypeKey]) > 0 {
		if resp.Header[ContentTypeKey][0] == ContentTypeGBK {
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
	isfirst := false
	for k, v := range params {
		if isfirst {
			s = s + "&"
		}
		s = s + k + "=" + url.QueryEscape(typeSwitcher(v))
		isfirst = true
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
