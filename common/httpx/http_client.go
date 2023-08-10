package httpx

import (
	"net"
	"net/http"
	"time"
)

var Client *http.Client

func getHttpClient() *http.Client {
	if Client == nil {
		dialer := &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
		Client = &http.Client{
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				DialContext:           dialer.DialContext,
				TLSHandshakeTimeout:   10 * time.Second,
				MaxIdleConns:          30,
				IdleConnTimeout:       90 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			}}
	}
	return Client
}
