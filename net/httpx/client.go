package httpx

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-xuan/quanx/os/errorx"
)

var client *Client

type Client struct {
	strategy int
	client   *http.Client
}

func (c *Client) Do(httpRequest *http.Request) (*Response, error) {
	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		return nil, errorx.Wrap(err, "client.Do error")
	}
	resp := &Response{
		code:    httpResponse.StatusCode,
		cookies: httpResponse.Cookies(),
	}
	defer httpResponse.Body.Close()
	var body []byte
	if body, err = io.ReadAll(httpResponse.Body); err != nil {
		return resp, errorx.Wrap(err, "http.Response.Body read error")
	}
	resp.body = body
	return resp, nil

}

// GetClient 获取http客户端
func GetClient(strategy ...ClientStrategy) *Client {
	if len(strategy) > 0 {
		return strategy[0].Client()
	} else {
		return newHttpClient()
	}
}

const (
	httpStrategyCode = iota + 1
	proxyStrategyCode
	httpsStrategyCode
	httpsProxyStrategyCode
)

func newHttpClient() *Client {
	if client == nil || client.strategy != httpStrategyCode {
		client = &Client{strategy: httpStrategyCode, client: &http.Client{Transport: newTransport()}}
	}
	return client
}

func newHttpsClient(crt string) *Client {
	if client == nil || client.strategy != httpsStrategyCode {
		client = &Client{strategy: httpsStrategyCode, client: &http.Client{Transport: newTransport(crt)}}
	}
	return client
}

func newHttpProxyClient(proxyUrl string) *Client {
	if client == nil || client.strategy != proxyStrategyCode {
		var transport = newTransport()
		if proxyURL, err := url.Parse(proxyUrl); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
		client = &Client{strategy: proxyStrategyCode, client: &http.Client{Transport: transport}}
	}
	return client
}

func newHttpsProxyClient(proxyUrl string, crt string) *Client {
	if client == nil || client.strategy != httpsProxyStrategyCode {
		var transport = newTransport(crt)
		if proxyURL, err := url.Parse(proxyUrl); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
		client = &Client{strategy: httpsProxyStrategyCode, client: &http.Client{Transport: transport}}
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
