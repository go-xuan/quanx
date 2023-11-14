package gatewayx

import (
	"errors"
	"github.com/go-xuan/quanx/utils/filex"
	"strings"
	"time"

	"github.com/go-xuan/quanx/public/nacosx"
	"github.com/go-xuan/quanx/utils/structx"
)

type App struct {
	Name   string `json:"name"`   // 服务名称
	Group  string `json:"group"`  // 服务分组
	Prefix string `json:"prefix"` // Api根路由
	Url    string `json:"url"`    // Api-Url
	Auth   bool   `json:"auth"`   // Api是否鉴权
	Exempt string `json:"exempt"` // Api免鉴权
}

var Apps []*App

// 获取微服务addr
func GetServerHttpURL(group, dataId, uri string) (url string, exempt bool, err error) {
	if data, ok := nacosx.GetNacosConfigMonitor().GetConfigData(group, dataId); ok && data.Changed {
		// 将当前最新的content数据同步到servers
		err = structx.ParseBytesToPointer(filex.Suffix(dataId), []byte(data.Content), &Apps)
		if err != nil {
			err = errors.New("微服务网关列表同步失败 ：" + err.Error())
			return
		}
		// 更新nacos监控中配置值
		data.Changed = false
		data.UpdateTime = time.Now().UnixMilli()
	}
	for _, server := range Apps {
		if MatchUrl(uri, server.Url) {
			if server.Exempt != "" {
				var Exempts = strings.Split(server.Exempt, ",")
				for _, exemptUrl := range Exempts {
					if strings.TrimSpace(exemptUrl) == uri {
						exempt = true
					}
				}
			}
			url, err = nacosx.SelectOneHealthyInstance(server.Name, server.Group)
			if err != nil {
				err = errors.New("微服务实例未注册 ：" + err.Error())
				return
			}
			url += server.Prefix
			return
		}
	}
	err = errors.New("微服务网关路由未配置")
	return
}

// URL匹配
func MatchUrl(uri, rule string) bool {
	if strings.HasSuffix(rule, `/**`) {
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
