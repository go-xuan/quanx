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

var _client *Client

// GetClient 获取httpx客户端
func GetClient(category ...ClientCategory) *Client {
	if len(category) > 0 && category[0] != nil {
		return category[0].Client()
	}
	return newHttpClient()
}

// 初始化httpx客户端
func newClient(strategy int, crt, proxy string) *Client {
	return &Client{
		strategy: strategy,
		client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: newTransport(crt, proxy),
		},
	}
}

// Client httpx客户端
type Client struct {
	strategy int          // 客户端类型
	client   *http.Client // http客户端
}

// HttpClient 获取http客户端
func (c *Client) HttpClient() *http.Client {
	return c.client
}

func (c *Client) Do(httpRequest *http.Request) (*Response, error) {
	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		return nil, errorx.Wrap(err, "do http request error")
	}
	resp := &Response{
		status:  httpResponse.StatusCode,
		cookies: httpResponse.Cookies(),
	}
	defer httpResponse.Body.Close()
	var body []byte
	if body, err = io.ReadAll(httpResponse.Body); err != nil {
		return resp, errorx.Wrap(err, "read http response body error")
	}
	resp.body = body
	return resp, nil

}

const (
	httpStrategyCode = iota + 1
	proxyStrategyCode
	httpsStrategyCode
	httpsProxyStrategyCode
)

func newHttpClient() *Client {
	if _client == nil || _client.strategy != httpStrategyCode {
		_client = newClient(httpStrategyCode, "", "")
	}
	return _client
}

func newHttpsClient(crt string) *Client {
	if _client == nil || _client.strategy != httpsStrategyCode {
		_client = newClient(httpsStrategyCode, crt, "")
	}
	return _client
}

func newHttpProxyClient(proxy string) *Client {
	if _client == nil || _client.strategy != proxyStrategyCode {
		_client = newClient(proxyStrategyCode, "", proxy)
	}
	return _client
}

func newHttpsProxyClient(crt, proxy string) *Client {
	if _client == nil || _client.strategy != httpsProxyStrategyCode {
		_client = newClient(httpsProxyStrategyCode, crt, proxy)
	}
	return _client
}

func newTransport(crt, proxy string) *http.Transport {
	dialer := &net.Dialer{
		Timeout:   time.Second * 10,
		KeepAlive: time.Second * 10,
	}
	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          30,
		IdleConnTimeout:       time.Second * 90,
		TLSHandshakeTimeout:   time.Second * 10,
		ExpectContinueTimeout: time.Second,
	}
	if crt != "" {
		if pem, err := os.ReadFile(crt); err == nil {
			pool := x509.NewCertPool()
			if ok := pool.AppendCertsFromPEM(pem); !ok {
				transport.TLSClientConfig = &tls.Config{
					ClientCAs:          pool,
					InsecureSkipVerify: true,
				}
			}
		} else {
			panic(err)
		}
	}
	if u, err := url.Parse(proxy); err == nil {
		transport.Proxy = http.ProxyURL(u)
	}
	return transport
}

func newTLSClientConfig(crt string) *tls.Config {
	if crt != "" {
		if pem, err := os.ReadFile(crt); err != nil {
			panic(err)
		} else {
			pool := x509.NewCertPool()
			if ok := pool.AppendCertsFromPEM(pem); !ok {
				return &tls.Config{
					ClientCAs:          pool,
					InsecureSkipVerify: true,
				}
			}
		}
	}
	return nil
}
