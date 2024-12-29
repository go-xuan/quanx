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
func GetClient(strategy ...ClientStrategy) *Client {
	if len(strategy) > 0 && strategy[0] != nil {
		return strategy[0].Client()
	}
	return newHttpClient()
}

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
		code:    httpResponse.StatusCode,
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
		_client = &Client{
			strategy: httpStrategyCode,
			client: &http.Client{
				Timeout:   time.Second * 10,
				Transport: newTransport(),
			},
		}
	}
	return _client
}

func newHttpsClient(crt string) *Client {
	if _client == nil || _client.strategy != httpsStrategyCode {
		_client = &Client{
			strategy: httpsStrategyCode,
			client: &http.Client{
				Timeout:   time.Second * 10,
				Transport: newTransport(crt),
			},
		}
	}
	return _client
}

func newHttpProxyClient(proxyUrl string) *Client {
	if _client == nil || _client.strategy != proxyStrategyCode {
		var transport = newTransport()
		if proxyURL, err := url.Parse(proxyUrl); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
		_client = &Client{
			strategy: proxyStrategyCode,
			client: &http.Client{
				Timeout:   time.Second * 10,
				Transport: transport,
			},
		}
	}
	return _client
}

func newHttpsProxyClient(proxyUrl string, crt string) *Client {
	if _client == nil || _client.strategy != httpsProxyStrategyCode {
		var transport = newTransport(crt)
		if proxyURL, err := url.Parse(proxyUrl); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
		_client = &Client{
			strategy: httpsProxyStrategyCode,
			client: &http.Client{
				Timeout:   time.Second * 10,
				Transport: transport,
			},
		}
	}
	return _client
}

func newTransport(crt ...string) *http.Transport {
	dialer := &net.Dialer{
		Timeout:   time.Second * 10,
		KeepAlive: time.Second * 10,
	}
	var transport = &http.Transport{
		TLSClientConfig:       newTLSClientConfig(crt...),
		DialContext:           dialer.DialContext,
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          30,
		IdleConnTimeout:       time.Second * 90,
		TLSHandshakeTimeout:   time.Second * 10,
		ExpectContinueTimeout: time.Second,
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
