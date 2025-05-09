package httpx

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/go-xuan/quanx/base/errorx"
)

var (
	_client *Client    // http客户端
	mu      sync.Mutex // 互斥锁
)

// GetClient 获取httpx客户端
func GetClient(uses ...ClientUse) *Client {
	if len(uses) > 0 {
		option := &ClientOption{}
		for _, use := range uses {
			use(option)
		}
		return option.Client()
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
	crt      string       // 证书存储路径
	proxy    string       // 代理服务器地址
	client   *http.Client // http客户端
}

// HttpClient 获取http客户端
func (c *Client) HttpClient() *http.Client {
	return c.client
}

func (c *Client) Do(request *http.Request) (*Response, error) {
	if c.client == nil {
		return nil, errorx.New("client is nil")
	}
	response, err := c.client.Do(request)
	if err != nil {
		return nil, errorx.Wrap(err, "do http request error")
	}
	resp := &Response{
		status:  response.StatusCode,
		cookies: response.Cookies(),
	}
	if resp.body, err = io.ReadAll(response.Body); err != nil {
		return resp, errorx.Wrap(err, "read http response body error")
	}
	_ = response.Body.Close()
	return resp, nil
}

const (
	httpStrategyCode = iota + 1
	proxyStrategyCode
	httpsStrategyCode
	httpsProxyStrategyCode
)

func newHttpClient() *Client {
	mu.Lock()
	defer mu.Unlock()
	if _client == nil || _client.strategy != httpStrategyCode {
		_client = newClient(httpStrategyCode, "", "")
	}
	return _client
}

func newHttpsClient(crt string) *Client {
	mu.Lock()
	defer mu.Unlock()
	if _client == nil || _client.strategy != httpsStrategyCode {
		_client = newClient(httpsStrategyCode, crt, "")
	}
	return _client
}

func newHttpProxyClient(proxy string) *Client {
	mu.Lock()
	defer mu.Unlock()
	if _client == nil || _client.strategy != proxyStrategyCode {
		_client = newClient(proxyStrategyCode, "", proxy)
	}
	return _client
}

func newHttpsProxyClient(crt, proxy string) *Client {
	mu.Lock()
	defer mu.Unlock()
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
			if !pool.AppendCertsFromPEM(pem) {
				transport.TLSClientConfig = &tls.Config{
					ClientCAs:          pool,
					InsecureSkipVerify: true,
				}
			}
		} else {
			panic(err)
		}
	}
	if proxy != "" {
		if u, err := url.Parse(proxy); err == nil {
			transport.Proxy = http.ProxyURL(u)
		} else {
			panic(err)
		}
	}
	return transport
}
