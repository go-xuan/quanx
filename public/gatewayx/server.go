package gatewayx

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/public/nacosx"
)

// 微服务配置
var servers []*Server

type Server struct {
	Name  string   `json:"name"`  // 服务名称
	Group string   `json:"group"` // 服务分组
	Url   string   `json:"url"`   // 微服务Url
	Auth  bool     `json:"auth"`  // 是否鉴权
	Skip  []string `json:"skip"`  // 跳过鉴权
}

// 获取微服务url
func GetServerHttpUrl(group, dataId, uri string) (string, bool, error) {
	configData, ok := nacosx.GetNacosConfigMonitor().GetConfigData(group, dataId)
	if ok && configData.Changed {
		err := json.Unmarshal([]byte(configData.Content), &servers)
		if err != nil {
			err = errors.New("微服务网关清单同步失败 ：" + err.Error())
			return "", false, err
		}
		// 更新nacos监控中配置值
		configData.Changed = false
		configData.UpdateTime = time.Now().UnixMilli()
	}
	for _, server := range servers {
		if MatchUrl(uri, server.Url) {
			var auth = server.Auth
			if server.Skip != nil && len(server.Skip) > 0 {
				for _, item := range server.Skip {
					if item == uri {
						auth = false
					}
				}
			}
			url, err := nacosx.SelectOneHealthyInstance(server.Name, server.Group)
			if err != nil {
				err = errors.New("微服务实例未注册 ：" + err.Error())
				return "", false, err
			}
			return constx.HttpPrefix + url, auth, nil
		}
	}
	return "", false, errors.New("微服务网关路由未配置")
}

// URL匹配
func MatchUrl(uri, rule string) bool {
	if strings.HasSuffix(rule, `/**`) {
		return strings.HasPrefix(uri, strings.TrimSuffix(rule, `/**`))
	} else if strings.HasSuffix(rule, `/*`) {
		prefix := strings.TrimSuffix(rule, `/*`)
		if strings.HasPrefix(uri, prefix) {
			uri = uri[len(prefix):]
			return !strings.Contains(uri, `/`)
		}
	}
	return false
}
