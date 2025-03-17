package enumx

type Enum[KT comparable, VT any] struct {
	keys []KT       // 保证有序
	data map[any]VT // 存储枚举值
}

// NewEnum 任意KV类型
func NewEnum[KT comparable, VT any]() *Enum[KT, VT] {
	return &Enum[KT, VT]{
		keys: make([]KT, 0),
		data: make(map[any]VT),
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
		e.data = make(map[any]VT)
	}
	e.keys = append(e.keys, k)
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

type StringEnum[T any] struct {
	keys []string     // 保证有序
	data map[string]T // 存储枚举值
}

// NewStringEnum K为string，v为任意类型
func NewStringEnum[T any]() *StringEnum[T] {
	return &StringEnum[T]{
		keys: make([]string, 0),
		data: make(map[string]T),
	}
}

func (e *StringEnum[T]) Get(k string) T {
	return e.data[k]
}

func (e *StringEnum[T]) Exist(k string) bool {
	if _, ok := e.data[k]; ok {
		return true
	}
	return false
}

func (e *StringEnum[T]) Add(k string, v T) *StringEnum[T] {
	if e.data == nil {
		e.keys = make([]string, 0)
		e.data = make(map[string]T)
	}
	e.keys = append(e.keys, k)
	e.data[k] = v
	return e
}

func (e *StringEnum[T]) Keys() []string {
	return e.keys
}

func (e *StringEnum[T]) Values() []T {
	var values []T
	for _, key := range e.keys {
		values = append(values, e.data[key])
	}
	return values
}

type IntEnum[T any] struct {
	keys []int     // 保证有序
	data map[int]T // 存储枚举值
}

func NewIntEnum[T any]() *IntEnum[T] {
	return &IntEnum[T]{
		keys: make([]int, 0),
		data: make(map[int]T),
	}
}

func (e *IntEnum[T]) Get(k int) T {
	return e.data[k]
}

func (e *IntEnum[T]) Exist(k int) bool {
	if _, ok := e.data[k]; ok {
		return true
	}
	return false
}

func (e *IntEnum[T]) Add(k int, v T) *IntEnum[T] {
	if e.data == nil {
		e.keys = make([]int, 0)
		e.data = make(map[int]T)
	}
	e.keys = append(e.keys, k)
	e.data[k] = v
	return e
}

func (e *IntEnum[T]) Keys() []int {
	return e.keys
}

func (e *IntEnum[T]) Values() []T {
	var values []T
	for _, key := range e.keys {
		values = append(values, e.data[key])
	}
	return values
}
