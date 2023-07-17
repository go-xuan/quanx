package gatewayx

import (
	"encoding/json"
	"errors"
	"github.com/quanxiaoxuan/quanx/common/httpx"
	"github.com/quanxiaoxuan/quanx/middleware/nacosx"
	"strings"
	"time"
)

var microServerList []*MicroServer

type MicroServer struct {
	Name  string   `json:"name"`  // 服务名称
	Group string   `json:"group"` // 服务分组
	Url   string   `json:"url"`   // 微服务Url
	Auth  bool     `json:"auth"`  // 是否鉴权
	Skip  []string `json:"skip"`  // 跳过鉴权
}

// 获取微服务url
func GetMicroServerHttpUrl(group, dataId, uri string) (string, bool, error) {
	configData, ok := nacosx.ConfigMonitor.GetConfigData(group, dataId)
	if ok && configData.Changed {
		err := json.Unmarshal([]byte(configData.Content), &microServerList)
		if err != nil {
			err = errors.New("微服务网关清单同步失败 ：" + err.Error())
			return "", false, err
		}
		// 更新nacos监控中配置值
		configData.Changed = false
		configData.UpdateTime = time.Now().UnixMilli()
	}
	for _, microServer := range microServerList {
		if MatchUrl(uri, microServer.Url) {
			var auth = microServer.Auth
			if microServer.Skip != nil && len(microServer.Skip) > 0 {
				for _, item := range microServer.Skip {
					if item == uri {
						auth = false
					}
				}
			}
			url, err := nacosx.SelectOneHealthyInstance(microServer.Name, microServer.Group)
			if err != nil {
				err = errors.New("微服务实例未注册 ：" + err.Error())
				return "", false, err
			}
			return httpx.HttpPrefix + url, auth, nil
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
