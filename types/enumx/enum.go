package enumx

type enum[KT any, VT any] struct {
	keys []KT       // 保证有序
	data map[any]VT // 存储枚举值
}

// NewEnum 任意KV类型
func NewEnum[KT any, VT any]() *enum[KT, VT] {
	return &enum[KT, VT]{
		keys: make([]KT, 0),
		data: make(map[any]VT),
	}
}

func (e *enum[KT, VT]) Get(k KT) VT {
	return e.data[k]
}

func (e *enum[KT, VT]) Add(k KT, v VT) *enum[KT, VT] {
	if e.data == nil {
		e.keys = make([]KT, 0)
		e.data = make(map[any]VT)
	}
	e.keys = append(e.keys, k)
	e.data[k] = v
	return e
}

func (e *enum[KT, VT]) Keys() []KT {
	return e.keys
}

func (e *enum[KT, VT]) Values() []VT {
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

type iIntEnum[T any] struct {
	keys []int     // 保证有序
	data map[int]T // 存储枚举值
}

func NewIntEnum[T any]() *iIntEnum[T] {
	return &iIntEnum[T]{
		keys: make([]int, 0),
		data: make(map[int]T),
	}
}

func (e *iIntEnum[T]) Get(k int) T {
	return e.data[k]
}

func (e *iIntEnum[T]) Add(k int, v T) *iIntEnum[T] {
	if e.data == nil {
		e.keys = make([]int, 0)
		e.data = make(map[int]T)
	}
	e.data = make(map[int]T)
	e.keys = append(e.keys, k)
	e.data[k] = v
	return e
}

func (e *iIntEnum[T]) Keys() []int {
	return e.keys
}

func (e *iIntEnum[T]) Values() []T {
	var values []T
	for _, key := range e.keys {
		values = append(values, e.data[key])
	}
	return values
}
