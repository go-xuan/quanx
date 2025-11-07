package gatewayx

import (
	"fmt"
	"strings"

	"github.com/go-xuan/quanx/serverx"
	"github.com/go-xuan/utilx/errorx"
)

// GetServerProxyURL 获取服务代理url
func GetServerProxyURL(url string) (string, error) {
	for _, proxy := range Gateway.Proxies {
		if proxy.Match(url) {
			if instance, err := serverx.SelectOne(proxy.Server); err != nil {
				return "", errorx.Wrap(err, "proxy instance not found")
			} else {
				return fmt.Sprintf("http://%s:%d/%s",
					instance.GetHost(),
					instance.GetPort(),
					strings.TrimPrefix(proxy.Prefix, "/")), nil
			}
		}
	}
	return "", errorx.Errorf("this url [%s] not match any gateway proxy", url)
}
