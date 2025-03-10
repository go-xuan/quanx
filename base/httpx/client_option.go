package httpx

// ClientOption http客户端设置
type ClientOption struct {
	crt   string // 证书存储路径
	proxy string // 代理服务器URL
}

func (c *ClientOption) Client() *Client {
	if c.crt == "" && c.proxy == "" {
		return newHttpClient()
	} else if c.crt != "" {
		return newHttpsClient(c.crt)
	} else if c.proxy != "" {
		return newHttpsClient(c.proxy)
	} else {
		return newHttpsProxyClient(c.crt, c.proxy)
	}
}

type ClientUse = func(opt *ClientOption)

func UseHttps(crt string) ClientUse {
	return func(opt *ClientOption) {
		opt.crt = crt
	}
}

func UseProxy(proxy string) ClientUse {
	return func(opt *ClientOption) {
		opt.proxy = proxy
	}
}
