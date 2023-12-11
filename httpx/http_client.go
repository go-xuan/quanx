package httpx

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"os"
	"time"
)

var client *http.Client

// 默认只支持http
func httpClient() *http.Client {
	if client == nil {
		client = NewClient()
	}
	return client
}

// 使用https必须先执行此初始化方法加载证书
func InitHttps(crt string) {
	if client == nil {
		client = NewClient(crt)
	}
}

func NewClient(crt ...string) *http.Client {
	dialer := &net.Dialer{Timeout: 3 * time.Second, KeepAlive: 10 * time.Second}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:       NewTLSClientConfig(crt...),
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			MaxIdleConns:          30,
			IdleConnTimeout:       90 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}}
}

func NewTLSClientConfig(crt ...string) *tls.Config {
	if len(crt) > 0 {
		pem, err := os.ReadFile(crt[0])
		if err != nil {
			panic(err)
		}
		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM(pem); !ok {
			panic(err)
		}
		return &tls.Config{
			ClientCAs:          pool,
			InsecureSkipVerify: true,
		}
	}
	return nil
}
