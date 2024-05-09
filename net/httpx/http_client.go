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
	Http       = "HTTP"
	Proxy      = "Proxy"
	Https      = "HTTPS"
	HttpsProxy = "httpsProxy"
	GET        = "GET"
	POST       = "POST"
	PUT        = "PUT"
	DELETE     = "DELETE"
)

var client *Client

type Client struct {
	mode   string
	client *http.Client
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// 默认只支持http
func httpClient() *Client {
	if client == nil {
		client = newHttpClient()
	}
	return client
}

// 切换客户端
func SwitchClient(modeAndParam ...string) *Client {
	if len(modeAndParam) == 2 && modeAndParam[0] == Https {
		client = newHttpsClient(modeAndParam[1])
	} else if len(modeAndParam) == 2 && modeAndParam[0] == Proxy {
		client = newHttpProxyClient(modeAndParam[1])
	} else if len(modeAndParam) == 3 && modeAndParam[0] == HttpsProxy {
		client = newHttpsProxyClient(modeAndParam[1], modeAndParam[2])
	} else {
		client = newHttpClient()
	}
	return client
}

func newHttpClient() *Client {
	if client == nil || client.mode != Http {
		client = &Client{mode: Http, client: &http.Client{Transport: newTransport()}}
	}
	return client
}

func newHttpsClient(crt string) *Client {
	if client == nil || client.mode != Https {
		client = &Client{mode: Https, client: &http.Client{Transport: newTransport(crt)}}
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
	if client == nil || client.mode != HttpsProxy {
		var transport = newTransport(crt)
		if proxyURL, err := url.Parse(proxyUrl); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
		client = &Client{mode: HttpsProxy, client: &http.Client{Transport: transport}}
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
