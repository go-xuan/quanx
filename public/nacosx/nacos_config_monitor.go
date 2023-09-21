package nacosx

import (
	"sync"
	"time"
)

var ConfigMonitor *Monitor

// nacos配置监听
type Monitor struct {
	mu        sync.RWMutex                       // 互斥锁
	Datas     map[string]map[string]*MonitorData // 配置数据
	ConfigNum int                                // 配置数量
}

// 监听配置数据
type MonitorData struct {
	Group      string // 配置分组
	DataId     string // 配置文件ID
	Content    string // 配置内容
	Changed    bool   // 修改标识
	UpdateTime int64  // 修改时间
}

// 初始化nacos配置监听
func GetNacosConfigMonitor() *Monitor {
	if ConfigMonitor == nil {
		ConfigMonitor = &Monitor{
			Datas:     make(map[string]map[string]*MonitorData),
			ConfigNum: 0,
		}
	}
	return ConfigMonitor
}

// 获取nacos配置
func (m *Monitor) GetConfigData(group, dataId string) (*MonitorData, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	gm, hasGroup := m.Datas[group]
	if hasGroup && gm != nil {
		data, hasData := gm[dataId]
		if hasData && data != nil {
			return data, true
		}
	}
	return nil, false
}

// 新增nacos配置监听
func (m *Monitor) AddConfigData(group, dataId, content string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var newData = MonitorData{group, dataId, content, false, time.Now().UnixMilli()}
	gm, hasGroup := m.Datas[group]
	if hasGroup && gm != nil {
		gm[dataId] = &newData
		m.Datas[group] = gm
	} else {
		var newGroup = make(map[string]*MonitorData)
		newGroup[dataId] = &newData
		m.Datas[group] = newGroup
	}
	m.ConfigNum++
	return
}

// 修改nacos配置信息
func (m *Monitor) UpdateConfigData(group, dataId, content string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	gm, hasGroup := m.Datas[group]
	if hasGroup && gm != nil {
		data, hasData := gm[dataId]
		if hasData && data != nil {
			data.Content = content
			data.Changed = true
			data.UpdateTime = time.Now().UnixMilli()
			return true
		}
	}
	return false
}
