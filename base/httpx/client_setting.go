package httpx

// ClientSetting http客户端设置
type ClientSetting struct {
	crt   string // 证书存储路径
	proxy string // 代理服务器URL
}

func (s *ClientSetting) Client() *Client {
	if s.crt == "" && s.proxy == "" {
		return newHttpClient()
	} else if s.crt != "" {
		return newHttpsClient(s.crt)
	} else if s.proxy != "" {
		return newHttpsClient(s.proxy)
	} else {
		return newHttpsProxyClient(s.crt, s.proxy)
	}
}

// Option 客户端选项
type Option = func(setting *ClientSetting)

func UseHttps(crt string) Option {
	return func(setting *ClientSetting) {
		setting.crt = crt
	}
}

func UseProxy(proxy string) Option {
	return func(setting *ClientSetting) {
		setting.proxy = proxy
	}
}
