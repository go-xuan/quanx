package idx

import (
	"sync"
)

var seqManager *SeqManager

type SeqManager struct {
	m map[string]*sequence
	l []*sequence
}

func Sequence() *SeqManager {
	if seqManager == nil {
		seqManager = &SeqManager{
			m: make(map[string]*sequence),
			l: make([]*sequence, 0),
		}
	}
	return seqManager
}

// Create 创建序列
func (m *SeqManager) Create(name string, start int, incr int) {
	var seq = &sequence{new(sync.RWMutex), name, start, incr, start}
	m.m[name] = seq
	m.l = append(m.l, seq)
}

// CurrVal 获取序列当前值
func (m *SeqManager) CurrVal(name string) int {
	if seq, ok := m.m[name]; ok {
		return seq.curr()
	} else {
		m.Create(name, 0, 1)
		return 0
	}
}

// NextVal 获取序列下值
func (m *SeqManager) NextVal(name string) int {
	if seq, ok := m.m[name]; ok {
		return seq.next()
	} else {
		m.Create(name, 1, 1)
		return 1
	}
}

// NextBatch 获取序列当前值
func (m *SeqManager) NextBatch(name string, n int) int {
	if seq, ok := m.m[name]; ok {
		var next = seq.next()
		seq.set(next + (n-1)*seq.increment)
		return next
	} else {
		m.Create(name, n+1, 1)
		return 1
	}
}

// Set 设置序列当前值
func (m *SeqManager) Set(name string, value int) {
	if seq, ok := m.m[name]; ok {
		seq.set(value)
	} else {
		m.Create(name, value, 1)
	}
}

// Reset 序列重置
func (m *SeqManager) Reset(name string) {
	if seq, ok := m.m[name]; ok {
		seq.reset()
	} else {
		m.Create(name, 0, 1)
	}
}

type sequence struct {
	*sync.RWMutex
	name      string // 序列名
	start     int    // 开始值
	increment int    // 递增值
	val       int    // 序列号
}

func (seq *sequence) curr() int {
	seq.RLock()
	defer seq.RUnlock()
	return seq.val
}

func (seq *sequence) next() int {
	seq.Lock()
	defer seq.Unlock()
	seq.val += seq.increment
	return seq.val
}

func (seq *sequence) reset() {
	seq.Lock()
	defer seq.Unlock()
	seq.val = seq.start
}

func (seq *sequence) set(v int) {
	seq.Lock()
	defer seq.Unlock()
	if v < seq.start {
		seq.val = seq.start
	} else {
		seq.val = v
	}
}
