package gatewayx

import (
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/stringx"

	"github.com/go-xuan/quanx/serverx"
)

// GetServerProxyURL 获取服务代理url
func GetServerProxyURL(url string) (string, error) {
	for _, proxy := range Gateway.Proxies {
		if proxy.Match(url) {
			if instance, err := serverx.GetCenter().SelectOne(proxy.Server); err != nil {
				return "", errorx.Wrap(err, "proxy instance not found")
			} else {
				return instance.GetDomain() + stringx.AddPrefix(proxy.Prefix, "/"), nil
			}
		}
	}
	return "", errorx.Errorf("this url [%s] not match any gateway proxy", url)
}
