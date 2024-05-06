package gatewayx

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-xuan/quanx/core/confx"
	"github.com/go-xuan/quanx/core/ginx"
	"github.com/go-xuan/quanx/core/nacosx"
	"github.com/go-xuan/quanx/types/stringx"
)

var Servers = &Gateway{}

type Gateway []*Server

func (g *Gateway) Theme() string {
	return "gateway"
}

func (g *Gateway) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "gateway.yaml",
		NacosDataId: "gateway.yaml",
		Listen:      true,
	}
}

func (g *Gateway) Run() error {
	return nil
}

// 微服务网关配置
type Server struct {
	Name   string   `yaml:"name" json:"name"`     // 微服务名称
	Group  string   `yaml:"group" json:"group"`   // 微服务分组
	Prefix string   `yaml:"prefix" json:"prefix"` // 微服务API前缀
	Router string   `yaml:"router" json:"router"` // API路由
	Auth   string   `yaml:"auth" json:"auth"`     // 鉴权方式(cookie/token/no)
	Ignore []string `yaml:"ignore" json:"ignore"` // 忽略鉴权(白名单)
}

// 获取微服务addr
func GetServerProxyAddr(group, dataId, url string) (addr string, auth string, err error) {
	if err = ListenConfigChanged(group, dataId); err != nil {
		err = errors.New("监听微服务网关配置失败 ：" + err.Error())
		return
	}
	for _, server := range *Servers {
		if MatchUrl(url, server.Router) {
			auth = server.Auth
			if auth != ginx.NoAuth && len(server.Ignore) > 0 {
				for _, item := range server.Ignore {
					if stringx.Index(url, strings.TrimSpace(item)) >= 0 {
						auth = ginx.NoAuth
						break
					}
				}
			}
			if addr, err = nacosx.SelectOneHealthyInstance(server.Name, server.Group); err != nil {
				err = errors.New("微服务实例未注册 ：" + err.Error())
				return
			}
			addr = "http://" + addr + server.Prefix
			return
		}
	}
	err = fmt.Errorf("未找到对应的网关路由配置，请检查微服务配置文件，或者确认请求接口[%s]是否正确", url)
	return

}

// 校验nacos配置
func ListenConfigChanged(group, dataId string) error {
	if data, ok := nacosx.GetNacosConfigMonitor().Get(group, dataId); ok {
		// 将当前最新的content数据同步到servers
		if err := data.Unmarshal(Servers); err != nil {
			return err
		}
		// 更新nacos监控中配置值
		data.SetChanged(false)
	}
	return nil
}

// URL匹配
func MatchUrl(uri, rule string) bool {
	if rule == "*" || rule == "/*" {
		return true
	} else if stringx.Index(rule, ",") >= 0 {
		var list = strings.Split(rule, ",")
		for _, item := range list {
			if stringx.Index(uri, strings.TrimSpace(item)) >= 0 {
				return true
			}
		}
	} else if strings.HasSuffix(rule, `/**`) {
		prefix := rule[:len(rule)-3]
		return strings.HasPrefix(uri, prefix)
	} else if strings.HasSuffix(rule, `/*`) {
		prefix := rule[:len(rule)-2]
		if strings.HasPrefix(uri, prefix) {
			uri = uri[len(prefix):]
			return stringx.Index(uri, `/`) < 0
		}
	} else {
		return stringx.Index(uri, rule) >= 0
	}
	return false
}
