package httpx

type ClientStrategy interface {
	Client() *Client
}

type HttpClient struct{}

func (c *HttpClient) Client() *Client {
	return newHttpClient()
}

type HttpsClient struct {
	Crt string
}

func (c *HttpsClient) Client() *Client {
	return newHttpsClient(c.Crt)
}

type ProxyClient struct {
	Proxy string
}

func (c *ProxyClient) Client() *Client {
	return newHttpProxyClient(c.Proxy)
}

type HttpsProxyClient struct {
	Proxy string
	Crt   string
}

func (c *HttpsProxyClient) Client() *Client {
	return newHttpsProxyClient(c.Proxy, c.Crt)
}
