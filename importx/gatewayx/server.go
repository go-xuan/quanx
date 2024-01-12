package gatewayx

import (
	"errors"
	"fmt"
	"github.com/go-xuan/quanx/importx/ginx"
	"strings"
	"time"

	"github.com/go-xuan/quanx/importx/nacosx"
	"github.com/go-xuan/quanx/utilx/filex"
	"github.com/go-xuan/quanx/utilx/marshalx"
)

var Apps []*App

// 微服务网关配置
type App struct {
	Name   string   `yaml:"name" json:"name"`     // 微服务名称
	Group  string   `yaml:"group" json:"group"`   // 微服务分组
	Prefix string   `yaml:"prefix" json:"prefix"` // 微服务API前缀
	Router string   `yaml:"router" json:"router"` // 路由规则
	Auth   string   `yaml:"auth" json:"auth"`     // 鉴权方式(cookie/token/no_auth)
	Skip   []string `yaml:"skip" json:"skip"`     // 跳过鉴权(白名单)
}

// 获取微服务addr
func GetServerProxyAddr(group, dataId, url string) (addr string, auth string, err error) {
	err = ListenConfigChanged(group, dataId)
	if err != nil {
		err = errors.New("监听微服务网关配置失败 ：" + err.Error())
		return
	}
	for _, server := range Apps {
		if MatchUrl(url, server.Router) {
			auth = server.Auth
			if auth != ginx.NoAuth && len(server.Skip) > 0 {
				for _, item := range server.Skip {
					if strings.Contains(url, strings.TrimSpace(item)) {
						auth = ginx.NoAuth
					}
				}
			}
			addr, err = nacosx.SelectOneHealthyInstance(server.Name, server.Group)
			if err != nil {
				err = errors.New("微服务实例未注册 ：" + err.Error())
				return
			}
			addr = "http://" + addr + server.Prefix
			return
		}
	}
	err = errors.New(fmt.Sprintf("未找到对应的网关路由配置，请检查微服务配置文件，或者确认请求接口[%s]是否正确", url))
	return

}

// 校验nacos配置
func ListenConfigChanged(group, dataId string) (err error) {
	data, ok := nacosx.GetNacosConfigMonitor().GetConfigData(group, dataId)
	if ok && data.Changed {
		// 将当前最新的content数据同步到servers
		err = marshalx.UnmarshalToPointer(&Apps, []byte(data.Content), filex.Suffix(dataId))
		if err != nil {
			return
		}
		// 更新nacos监控中配置值
		data.Changed = false
		data.UpdateTime = time.Now().UnixMilli()
	}
	return
}

// URL匹配
func MatchUrl(uri, rule string) bool {
	if rule == "*" || rule == "/*" {
		return true
	} else if strings.Contains(rule, ",") {
		var list = strings.Split(rule, ",")
		for _, item := range list {
			if strings.Contains(uri, strings.TrimSpace(item)) {
				return true
			}
		}
	} else if strings.HasSuffix(rule, `/**`) {
		return strings.HasPrefix(uri, strings.TrimSuffix(rule, `/**`))
	} else if strings.HasSuffix(rule, `/*`) {
		var prefix = strings.TrimSuffix(rule, `/*`)
		if strings.HasPrefix(uri, prefix) {
			uri = uri[len(prefix):]
			return !strings.Contains(uri, `/`)
		}
	} else {
		return strings.Contains(uri, rule)
	}
	return false
}
