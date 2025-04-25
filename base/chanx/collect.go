package chanx

import (
	"sync"

	"github.com/go-xuan/quanx/base/errorx"
	log "github.com/sirupsen/logrus"
)

func NewCollectManager[T any]() *CollectManager[T] {
	return &CollectManager[T]{
		mu:         &sync.Mutex{},
		wg:         &sync.WaitGroup{},
		collectors: make(map[string]*Collector[T]),
	}
}

type CollectManager[T any] struct {
	mu         *sync.Mutex
	wg         *sync.WaitGroup
	collectors map[string]*Collector[T]
}

func (m *CollectManager[T]) Add(key string, collector *Collector[T]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wg.Add(1)
	m.collectors[key] = collector
}

func (m *CollectManager[T]) Collect() {
	m.mu.Lock()
	// 复制一份 collectors 避免在 goroutine 执行期间修改 map 导致并发问题
	collectorsCopy := make(map[string]*Collector[T], len(m.collectors))
	for k, v := range m.collectors {
		collectorsCopy[k] = v
	}
	defer m.mu.Unlock()
	for _, collector := range collectorsCopy {
		go collector.Collect(m.wg)
	}
}

func (m *CollectManager[T]) Result(key string) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	collector, exists := m.collectors[key]
	if !exists {
		var zero T
		return zero, errorx.Errorf("collector with key %s not found", key)
	}
	return collector.Result()
}

func (m *CollectManager[T]) Close() {
	m.wg.Wait()
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, collector := range m.collectors {
		collector.Close()
	}
}

func NewCollector[T any](f func() (T, error)) *Collector[T] {
	return &Collector[T]{
		ch:    make(chan T, 1),
		errCh: make(chan error, 1),
		fn:    f,
	}
}

// Collector channel收集器
type Collector[T any] struct {
	ch    chan T            // channel管道
	errCh chan error        // 错误通道
	fn    func() (T, error) // 收集器方法
}

// Collect 收集数据
func (c *Collector[T]) Collect(wg *sync.WaitGroup) {
	defer wg.Done()
	if result, err := c.fn(); err == nil {
		c.ch <- result
	} else {
		c.errCh <- err
		log.WithError(err).Error("collector chan result error")
	}
}

func (c *Collector[T]) Result() (T, error) {
	var zero T
	select {
	case result := <-c.ch:
		return result, nil
	case err := <-c.errCh:
		return zero, err
	default:
		return zero, errorx.New("no result or error available yet")
	}
}

func (c *Collector[T]) Close() {
	close(c.ch)
	close(c.errCh)
}
