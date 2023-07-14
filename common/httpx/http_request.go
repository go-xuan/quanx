package httpx

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func GetHttp(url string) (result []byte, err error) {
	if !strings.HasPrefix(url, HttpPrefix) {
		url = HttpPrefix + url
	}
	var resp *http.Response
	resp, err = getHttpClient().Get(url)
	if err != nil {
		return
	}
	return io.ReadAll(resp.Body)
}

func PostHttp(url string, param interface{}) (result []byte, err error) {
	if !strings.HasPrefix(url, HttpPrefix) {
		url = HttpPrefix + url
	}
	var jsons []byte
	jsons, err = json.Marshal(param)
	if err != nil {
		return
	}
	var resp *http.Response
	resp, err = getHttpClient().Post(url, ContentTypeJson, bytes.NewBuffer(jsons))
	if err != nil {
		return
	}
	result, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

func PutHttp(url string, param interface{}) (result []byte, err error) {
	if !strings.HasPrefix(url, HttpPrefix) {
		url = HttpPrefix + url
	}
	var jsons []byte
	jsons, err = json.Marshal(param)
	if err != nil {
		return
	}
	var req *http.Request
	req, err = http.NewRequest(Put, url, bytes.NewBuffer(jsons))
	req.Header.Set(ContentTypeKey, ContentTypeJson)
	if err != nil {
		return
	}
	var resp *http.Response
	resp, err = getHttpClient().Do(req)
	if err != nil {
		return
	}
	return io.ReadAll(resp.Body)
}

func DeleteHttp(url string) (res []byte, err error) {
	if !strings.HasPrefix(url, HttpPrefix) {
		url = HttpPrefix + url
	}
	var req *http.Request
	req, err = http.NewRequest(Delete, url, nil)
	if err != nil {
		return
	}
	var resp *http.Response
	resp, err = getHttpClient().Do(req)
	if err != nil {
		return
	}
	return io.ReadAll(resp.Body)
}
