package nacosx

import (
	"sync"
	"time"
)

var monitor *Monitor

// nacos配置监听
type Monitor struct {
	mu        sync.RWMutex                       // 互斥锁
	Data      map[string]map[string]*MonitorData // 监听数据
	ConfigNum int                                // 监听数量
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
	if monitor == nil {
		monitor = &Monitor{
			Data:      make(map[string]map[string]*MonitorData),
			ConfigNum: 0,
		}
	}
	return monitor
}

// 获取nacos配置
func (m *Monitor) GetConfigData(group, dataId string) (data *MonitorData, exist bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var gm, ok = m.Data[group]
	if ok && gm != nil {
		data, exist = gm[dataId]
		return
	}
	return
}

// 新增nacos配置监听
func (m *Monitor) AddConfigData(group, dataId, content string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var newData = MonitorData{group, dataId, content, false, time.Now().UnixMilli()}
	groupData, hasGroup := m.Data[group]
	if hasGroup && groupData != nil {
		groupData[dataId] = &newData
		m.Data[group] = groupData
	} else {
		var newGroup = make(map[string]*MonitorData)
		newGroup[dataId] = &newData
		m.Data[group] = newGroup
	}
	m.ConfigNum++
	return
}

// 修改nacos配置信息
func (m *Monitor) UpdateConfigData(group, dataId, content string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	groupData, hasGroup := m.Data[group]
	if hasGroup && groupData != nil {
		data, hasData := groupData[dataId]
		if hasData && data != nil {
			data.Content = content
			data.Changed = true
			data.UpdateTime = time.Now().UnixMilli()
			return true
		}
	}
	return false
}
