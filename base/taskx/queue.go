package taskx

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

func Queue() *QueueScheduler {
	return &QueueScheduler{
		mutex: new(sync.Mutex),
		tasks: make(map[string]*QueueTask),
	}
}

// QueueScheduler 队列任务处理调度器
type QueueScheduler struct {
	mutex *sync.Mutex           // 锁
	head  *QueueTask            // 头部任务
	tail  *QueueTask            // 尾部任务
	tasks map[string]*QueueTask // 任务列表
}

func (q *QueueScheduler) Names() []string {
	var names = make([]string, 0, len(q.tasks))
	for t := q.head; t != nil; t = t.next {
		names = append(names, t.name)
	}
	return names
}

// Execute 执行队列任务
func (q *QueueScheduler) Execute() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for q.head != nil {
		logger := log.WithField("current_task", q.head.name)
		if q.head.next != nil {
			logger = logger.WithField("next_task", q.head.next.name)
		}
		logger.Info("queue task execute start")
		if err := q.head.fn(); err != nil {
			logger.Error("queue task execute error")
			return err
		}
		delete(q.tasks, q.head.name)
		q.head = q.head.next
	}
	if len(q.tasks) == 0 {
		q.head = nil
		q.tail = nil
	}

	return nil
}

// Add 新增队列任务（默认尾插）
func (q *QueueScheduler) Add(fn func() error, name string) {
	q.AddTail(fn, name)
}

// AddTail 尾插（当前新增任务添加到队列末尾）
func (q *QueueScheduler) AddTail(fn func() error, name string) {
	if name != "" && fn != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.tasks[name]; !exist {
			var task = &QueueTask{name: name, fn: fn}
			if tail := q.tail; tail != nil {
				task.prev = tail
				tail.next = task
				q.tail = task
			} else {
				q.head = task
				q.tail = task
			}
			q.tasks[name] = task
		} else {
			q.tasks[name].fn = fn
		}
	}
}

// AddHead 头插（当前新增任务添加到队列首位）
func (q *QueueScheduler) AddHead(fn func() error, name string) {
	if name != "" && fn != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.tasks[name]; !exist {
			var task = &QueueTask{name: name, fn: fn}
			if head := q.head; head != nil {
				head.prev = task
				task.next = head
				q.head = task
			} else {
				q.head = task
				q.tail = task
			}
			q.tasks[name] = task
		} else {
			q.tasks[name].fn = fn
		}
	}
}

// AddAfter 后插队(当前新增任务添加到after任务之后)
func (q *QueueScheduler) AddAfter(fn func() error, name, after string) {
	if name != "" && fn != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.tasks[name]; !exist {
			var task = &QueueTask{name: name, fn: fn}
			if target, ok := q.tasks[after]; ok && target.next != nil {
				task.next = target.next
				task.prev = target
				target.next.prev = task
				target.next = task
			} else {
				target = q.tail
				target.next = task
				task.prev = target
				q.tail = task
			}
			q.tasks[name] = task
		} else {
			q.tasks[name].fn = fn
		}
	}
}

// AddBefore 前插队(当前新增任务添加到before任务之后)
func (q *QueueScheduler) AddBefore(fn func() error, name string, before string) {
	if name != "" && fn != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.tasks[name]; !exist {
			var task = &QueueTask{name: name, fn: fn}
			if target, ok := q.tasks[before]; ok && target.prev != nil {
				task.prev = target.prev
				task.next = target
				target.prev.next = task
				target.prev = task
			} else {
				target = q.head
				target.prev = task
				task.next = target
				q.head = task
			}
			q.tasks[name] = task
		} else {
			q.tasks[name].fn = fn
		}
	}
}

func (q *QueueScheduler) Exist(name string) bool {
	if _, ok := q.tasks[name]; ok {
		return true
	}
	return false
}

// Remove 移除任务
func (q *QueueScheduler) Remove(name string) {
	if name != "" {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if task, ok := q.tasks[name]; ok {
			if task.prev == nil {
				task.next.prev = nil
				q.head = task.next
			} else if task.next == nil {
				task.prev.next = nil
				q.tail = task.prev
			} else {
				task.prev.next = task.next
				task.next.prev = task.prev
			}
			delete(q.tasks, name)
		}
	}
}

// Clear 清除所有队列任务
func (q *QueueScheduler) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.head = nil
	q.tail = nil
	q.tasks = make(map[string]*QueueTask)
}

// QueueTask 队列任务
type QueueTask struct {
	name string       // 任务名
	fn   func() error // 当前任务
	prev *QueueTask   // 指向上一个任务
	next *QueueTask   // 指向下一个任务
}

func (t *QueueTask) HasNext() bool {
	return t.next != nil
}

func (t *QueueTask) Name() string {
	return t.name
}
