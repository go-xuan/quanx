package httpx

// ClientCategory 客户端类型
type ClientCategory interface {
	Client() *Client
}

// ClientConfig 客户端通用类型
type ClientConfig struct {
	Crt   string // 证书存储路径
	Proxy string // 代理服务器URL
}

func (c *ClientConfig) Client() *Client {
	if c.Crt == "" && c.Proxy == "" {
		return newHttpClient()
	} else if c.Crt != "" {
		return newHttpsClient(c.Crt)
	} else if c.Proxy != "" {
		return newHttpsClient(c.Proxy)
	} else {
		return newHttpsProxyClient(c.Crt, c.Proxy)
	}
}

// HttpClientCategory 默认http客户端
type HttpClientCategory struct{}

func (c *HttpClientCategory) Client() *Client {
	return newHttpClient()
}

// HttpsClientCategory https客户端
type HttpsClientCategory struct {
	Crt string // 证书存储路径
}

func (c *HttpsClientCategory) Client() *Client {
	return newHttpsClient(c.Crt)
}

// ProxyClientCategory 代理客户端
type ProxyClientCategory struct {
	Proxy string // 代理服务器URL
}

func (c *ProxyClientCategory) Client() *Client {
	return newHttpProxyClient(c.Proxy)
}

// HttpsProxyClientCategory https代理客户端
type HttpsProxyClientCategory struct {
	Crt   string // 证书存储路径
	Proxy string // 代理服务器URL
}

func (c *HttpsProxyClientCategory) Client() *Client {
	return newHttpsProxyClient(c.Crt, c.Proxy)
}
