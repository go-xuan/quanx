package enumx

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
	keys []KT      // 保证有序
	data map[KT]VT // 存储枚举值
}

func (e *Enum[KT, VT]) Len() int {
	return len(e.keys)
}

func (e *Enum[KT, VT]) Clear() {
	e.keys = make([]KT, 0)
	e.data = make(map[KT]VT)
}

func (e *Enum[KT, VT]) Remove(k KT) {
	delete(e.data, k)
	for i, key := range e.keys {
		if key == k {
			e.keys = append(e.keys[:i], e.keys[i+1:]...)
		}
	}
}

func (e *Enum[KT, VT]) Get(k KT) VT {
	return e.data[k]
}

func (e *Enum[KT, VT]) Exist(k KT) bool {
	if _, ok := e.data[k]; ok {
		return true
	}
	return false
}

func (e *Enum[KT, VT]) Add(k KT, v VT) *Enum[KT, VT] {
	if e.data == nil {
		e.keys = make([]KT, 0)
		e.data = make(map[KT]VT)
	}
	if _, ok := e.data[k]; !ok {
		e.keys = append(e.keys, k)
	}
	e.data[k] = v
	return e
}

func (e *Enum[KT, VT]) Keys() []KT {
	return e.keys
}

func (e *Enum[KT, VT]) Values() []VT {
	var values []VT
	for _, key := range e.keys {
		values = append(values, e.data[key])
	}
	return values
}
