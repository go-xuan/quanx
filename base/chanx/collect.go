package chanx

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

func NewCollectManager[T any]() *CollectManager[T] {
	return &CollectManager[T]{
		wg:         &sync.WaitGroup{},
		collectors: make(map[string]*Collector[T]),
	}
}

type CollectManager[T any] struct {
	wg         *sync.WaitGroup
	collectors map[string]*Collector[T]
}

func (m *CollectManager[T]) Add(key string, collector *Collector[T]) {
	m.wg.Add(1)
	m.collectors[key] = collector
}

func (m *CollectManager[T]) Collect() {
	for _, collector := range m.collectors {
		go collector.Collect(m.wg)
	}
}

func (m *CollectManager[T]) Result(key string, result any) {
	result = m.collectors[key].Result()
}

func (m *CollectManager[T]) Close() {
	m.wg.Wait()
	for _, collector := range m.collectors {
		collector.Close()
	}
}

func NewCollector[T any](f func() (T, error)) *Collector[T] {
	return &Collector[T]{
		ch: make(chan T, 1),
		fn: f,
	}
}

// Collector channel收集器
type Collector[T any] struct {
	ch chan T            // channel管道
	fn func() (T, error) // 收集器方法
}

// Collect 收集数据
func (c *Collector[T]) Collect(wg *sync.WaitGroup) {
	defer wg.Done()
	if result, err := c.fn(); err == nil {
		c.ch <- result
	} else {
		log.WithError(err).Error("collector chan result error")
	}
}

func (c *Collector[T]) Result() T {
	return <-c.ch
}

func (c *Collector[T]) Close() {
	close(c.ch)
}
