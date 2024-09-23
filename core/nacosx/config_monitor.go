package nacosx

import (
	"sync"
	"time"

	"github.com/go-xuan/quanx/utils/marshalx"
)

var monitor *ConfigMonitor

// GetConfigMonitor 初始化nacos配置监听
func GetConfigMonitor() *ConfigMonitor {
	if monitor == nil {
		monitor = &ConfigMonitor{
			data: make(map[string]*ConfigData),
			num:  0,
		}
	}
	return monitor
}

// ConfigMonitor nacos配置监听器
type ConfigMonitor struct {
	mu   sync.RWMutex           // 互斥锁
	data map[string]*ConfigData // 配置数据
	num  int                    // 配置数量
}

// ConfigData 监听配置数据
type ConfigData struct {
	group   string // 配置分组
	dataId  string // 配置DataId
	content string // 配置内容
	changed bool   // 已修改标识
	modify  int64  // 最近修改时间
}

func (d *ConfigData) SetChanged(changed bool) {
	d.changed = changed
	d.modify = time.Now().UnixMilli()
}

func (d *ConfigData) Unmarshal(v any) error {
	return marshalx.NewCase(d.dataId).Unmarshal([]byte(d.content), v)
}

func getKey(group, dataId string) string {
	return group + "_" + dataId
}

// Set 新增nacos配置监听
func (m *ConfigMonitor) Set(group, dataId, content string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var key = getKey(group, dataId)
	if data, exist := m.data[key]; !exist {
		m.data[key] = &ConfigData{
			group,
			dataId,
			content,
			false,
			time.Now().UnixMilli()}
		m.num++
	} else {
		data.content = content
		data.SetChanged(true)
	}
	return
}

// Get 获取nacos配置
func (m *ConfigMonitor) Get(group, dataId string) (data *ConfigData, exist bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, exist = m.data[getKey(group, dataId)]
	return
}
