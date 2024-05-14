package idx

import (
	"github.com/go-xuan/quanx/types/anyx"
	"sync"
)

var seqManager *SequenceManager

type SequenceManager struct {
	*sync.RWMutex
	seqMap  map[string]*SequenceData
	seqList []*SequenceData
}

func Sequence() *SequenceManager {
	if seqManager == nil {
		seqManager = &SequenceManager{
			RWMutex: new(sync.RWMutex),
			seqList: make([]*SequenceData, 0),
			seqMap:  make(map[string]*SequenceData),
		}
	}
	return seqManager
}

func (m *SequenceManager) All() []*SequenceData {
	m.RLock()
	defer m.RUnlock()
	return m.seqList
}

func (m *SequenceManager) Add(name string, start int64, incr ...int64) {
	m.Lock()
	defer m.Unlock()
	var seq = &SequenceData{new(sync.RWMutex), name, start, anyx.Default(1, incr...), start}
	m.seqMap[name] = seq
	m.seqList = append(m.seqList, seq)
}

// 获取序列当前值
func (m *SequenceManager) CurrVal(name string) int64 {
	m.RLock()
	defer m.RUnlock()
	if seq, ok := m.seqMap[name]; ok {
		return seq.Curr()
	}
	return -1
}

// 获取序列下个值
func (m *SequenceManager) NextVal(name string) int64 {
	m.RLock()
	defer m.RUnlock()
	if seq, ok := m.seqMap[name]; ok {
		return seq.Next()
	}
	return -1
}

// 设置序列当前值
func (m *SequenceManager) SetVal(name string, value int64) {
	m.RLock()
	defer m.RUnlock()
	if seq, ok := m.seqMap[name]; ok {
		seq.set(value)
	}
}

// 序列重置
func (m *SequenceManager) ResetVal(name string) {
	m.RLock()
	defer m.RUnlock()
	if seq, ok := m.seqMap[name]; ok {
		seq.Reset()
	}
}

type SequenceData struct {
	*sync.RWMutex
	name      string // 序列名
	start     int64  // 开始值
	increment int64  // 递增值
	val       int64  // 序列号
}

func (seq *SequenceData) Curr() int64 {
	seq.RLock()
	defer seq.RUnlock()
	return seq.val
}

func (seq *SequenceData) Next() int64 {
	seq.Lock()
	defer seq.Unlock()
	seq.val += seq.increment
	return seq.val
}

func (seq *SequenceData) Reset() {
	seq.Lock()
	defer seq.Unlock()
	seq.val = seq.start
}

func (seq *SequenceData) set(v int64) {
	seq.Lock()
	defer seq.Unlock()
	if v < seq.start {
		seq.val = seq.start
	} else {
		seq.val = v
	}
}
