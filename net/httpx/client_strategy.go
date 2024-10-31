package httpx

// ClientStrategy 客户端类型
type ClientStrategy interface {
	Client() *Client
}

type HttpClientStrategy struct{}

func (c *HttpClientStrategy) Client() *Client {
	return newHttpClient()
}

type HttpsClientStrategy struct {
	Crt string
}

func (c *HttpsClientStrategy) Client() *Client {
	return newHttpsClient(c.Crt)
}

type ProxyClientStrategy struct {
	Proxy string
}

func (c *ProxyClientStrategy) Client() *Client {
	return newHttpProxyClient(c.Proxy)
}

type HttpsProxyClientStrategy struct {
	Proxy string
	Crt   string
}

func (c *HttpsProxyClientStrategy) Client() *Client {
	return newHttpsProxyClient(c.Proxy, c.Crt)
}
