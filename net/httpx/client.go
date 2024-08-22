package httpx

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	Http = iota + 1
	Proxy
	https
	httpsProxy
)

var client *Client

type Client struct {
	mode   int
	client *http.Client
}

// GetClient 获取http客户端
func GetClient(strategy ...ClientStrategy) *Client {
	if len(strategy) > 0 {
		return strategy[0].Client()
	} else {
		return newHttpClient()
	}
}

func newHttpClient() *Client {
	if client == nil || client.mode != Http {
		client = &Client{mode: Http, client: &http.Client{Transport: newTransport()}}
	}
	return client
}

func newHttpsClient(crt string) *Client {
	if client == nil || client.mode != https {
		client = &Client{mode: https, client: &http.Client{Transport: newTransport(crt)}}
	}
	return client
}

func newHttpProxyClient(proxyUrl string) *Client {
	if client == nil || client.mode != Proxy {
		var transport = newTransport()
		if proxyURL, err := url.Parse(proxyUrl); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
		client = &Client{mode: Proxy, client: &http.Client{Transport: transport}}
	}
	return client
}

func newHttpsProxyClient(proxyUrl string, crt string) *Client {
	if client == nil || client.mode != httpsProxy {
		var transport = newTransport(crt)
		if proxyURL, err := url.Parse(proxyUrl); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
		client = &Client{mode: httpsProxy, client: &http.Client{Transport: transport}}
	}
	return client
}

func newTransport(crt ...string) *http.Transport {
	dialer := &net.Dialer{Timeout: 3 * time.Second, KeepAlive: 10 * time.Second}
	var transport = &http.Transport{
		TLSClientConfig:       newTLSClientConfig(crt...),
		DialContext:           dialer.DialContext,
		Proxy:                 http.ProxyFromEnvironment,
		TLSHandshakeTimeout:   10 * time.Second,
		MaxIdleConns:          30,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return transport
}

func newTLSClientConfig(crt ...string) *tls.Config {
	if len(crt) > 0 && crt[0] != "" {
		if pem, err := os.ReadFile(crt[0]); err != nil {
			panic(err)
		} else {
			pool := x509.NewCertPool()
			if ok := pool.AppendCertsFromPEM(pem); !ok {
				panic(err)
			}
			return &tls.Config{
				ClientCAs:          pool,
				InsecureSkipVerify: true,
			}
		}
	}
	return nil
}
