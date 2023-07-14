package httpx

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var syncOnce sync.Once

func getHttpClient() *http.Client {
	var Client *http.Client
	syncOnce.Do(func() {
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
	})
	return Client
}
