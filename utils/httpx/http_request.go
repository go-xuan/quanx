package httpx

import (
	"bytes"
	"encoding/json"
	"github.com/quanxiaoxuan/quanx/common/constx"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func HttpGet(url string) (result []byte, err error) {
	var resp *http.Response
	resp, err = getHttpClient().Get(url)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	return io.ReadAll(resp.Body)
}

func HttpGetWithHeader(url string, header map[string]string) (result []byte, err error) {
	var req *http.Request
	req, err = http.NewRequest(constx.GET, url, nil)
	if err != nil {
		return
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	var resp *http.Response
	resp, err = getHttpClient().Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	return io.ReadAll(resp.Body)
}

func HttpPost(url string, body interface{}) (result []byte, err error) {
	bodyBytes, _ := json.Marshal(body)
	var resp *http.Response
	resp, err = getHttpClient().Post(url, constx.JsonContent, bytes.NewBuffer(bodyBytes))
	defer resp.Body.Close()
	if err != nil {
		return
	}
	return io.ReadAll(resp.Body)
}

func HttpPostWithHeader(url string, body interface{}, header map[string]string) (result []byte, err error) {
	bodyBytes, _ := json.Marshal(body)
	var req *http.Request
	req, err = http.NewRequest(constx.POST, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return
	}
	req.Header.Set(constx.ContentType, constx.JsonContent)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	var resp *http.Response
	resp, err = getHttpClient().Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	return io.ReadAll(resp.Body)
}

func HttpForm(url string, form url.Values) (result []byte, err error) {
	var resp *http.Response
	client := getHttpClient()
	resp, err = client.PostForm(url, form)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	return io.ReadAll(resp.Body)
}

func HttpFormWithHeader(url string, form url.Values, header map[string]string) (result []byte, err error) {
	var req *http.Request
	req, err = http.NewRequest(constx.POST, url, strings.NewReader(form.Encode()))
	if err != nil {
		return
	}
	req.Header.Set(constx.ContentType, constx.FormContent)
	for k, v := range header {
		req.Header.Set(k, v)
	}

	var resp *http.Response
	resp, err = getHttpClient().Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	return io.ReadAll(resp.Body)
}

func HttpPut(url string, param interface{}) (result []byte, err error) {
	var jsons []byte
	jsons, err = json.Marshal(param)
	if err != nil {
		return
	}
	var req *http.Request
	req, err = http.NewRequest(constx.PUT, url, bytes.NewBuffer(jsons))
	if err != nil {
		return
	}
	req.Header.Set(constx.ContentType, constx.JsonContent)
	var resp *http.Response
	resp, err = getHttpClient().Do(req)
	if err != nil {
		return
	}
	defer func(Body *http.Response) {
		_ = resp.Body.Close()
	}(resp)
	return io.ReadAll(resp.Body)
}

func HttpDelete(url string) (res []byte, err error) {
	var req *http.Request
	req, err = http.NewRequest(constx.DELETE, url, nil)
	if err != nil {
		return
	}
	var resp *http.Response
	resp, err = getHttpClient().Do(req)
	if err != nil {
		return
	}
	defer func(Body *http.Response) {
		_ = resp.Body.Close()
	}(resp)
	return io.ReadAll(resp.Body)
}
