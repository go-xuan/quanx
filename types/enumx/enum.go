package enumx

import "sync"

// NewStringEnum key为string类型，value为任意类型
func NewStringEnum[T any]() *Enum[string, T] {
	return &Enum[string, T]{
		keys: make([]string, 0),
		data: make(map[string]T),
	}
}

// NewIntEnum key为int类型，value为任意类型
func NewIntEnum[T any]() *Enum[int, T] {
	return &Enum[int, T]{
		keys: make([]int, 0),
		data: make(map[int]T),
	}
}

// NewEnum key为comparable类型，value为任意类型
func NewEnum[KT comparable, VT any]() *Enum[KT, VT] {
	return &Enum[KT, VT]{
		keys: make([]KT, 0),
		data: make(map[KT]VT),
	}
}

type Enum[KT comparable, VT any] struct {
	mu   sync.RWMutex // 读写锁
	keys []KT         // 保证有序
	data map[KT]VT    // 存储枚举值
}

func (e *Enum[KT, VT]) Len() int {
	return len(e.keys)
}

func (e *Enum[KT, VT]) Clear() {
	e.keys = make([]KT, 0)
	e.data = make(map[KT]VT)
}

func (e *Enum[KT, VT]) Remove(k KT) {
	if len(e.keys) == 0 {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.data, k)
	i := 0 // 使用双指针法删除切片元素
	for _, key := range e.keys {
		if key != k {
			e.keys[i] = key
			i++
		}
	}
	e.keys = e.keys[:i]
}

func (e *Enum[KT, VT]) Get(k KT) VT {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.data[k]
}

func (e *Enum[KT, VT]) Exist(k KT) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, ok := e.data[k]
	return ok
}

func (e *Enum[KT, VT]) Add(k KT, v VT) *Enum[KT, VT] {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.data[k]; !ok {
		e.keys = append(e.keys, k)
	}
	e.data[k] = v
	return e
}

func (e *Enum[KT, VT]) Keys() []KT {
	e.mu.RLock()
	defer e.mu.RUnlock()
	// 返回切片的副本，避免外部修改内部数据
	keysCopy := make([]KT, len(e.keys))
	copy(keysCopy, e.keys)
	return keysCopy
}

func (e *Enum[KT, VT]) Values() []VT {
	e.mu.Lock()
	defer e.mu.Unlock()
	var values []VT
	for _, key := range e.keys {
		values = append(values, e.data[key])
	}
	return values
}
