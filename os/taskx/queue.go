package taskx

import (
	"sync"
)

// QueueScheduler 队列任务处理调度器
type QueueScheduler struct {
	mutex *sync.Mutex           // 锁
	head  *QueueTask            // 首个任务
	tail  *QueueTask            // 末尾任务
	tasks map[string]*QueueTask // 所有任务
}

// QueueTask 队列任务
type QueueTask struct {
	name string     // 任务名
	fn   func()     // 当前任务
	prev *QueueTask // 指向上一个任务
	next *QueueTask // 指向下一个任务
}

func (t *QueueTask) HasNext() bool {
	return t.next != nil
}

func (t *QueueTask) Name() string {
	return t.name
}

func Queue() *QueueScheduler {
	return &QueueScheduler{
		mutex: new(sync.Mutex),
		tasks: make(map[string]*QueueTask),
	}
}

// Execute 执行队列任务
func (q *QueueScheduler) Execute() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for q.head != nil {
		q.head.fn()
		delete(q.tasks, q.head.name)
		q.head = q.head.next
	}
	if len(q.tasks) == 0 {
		q.head = nil
		q.tail = nil
	}
}

// Add 新增队列任务（默认尾插）
func (q *QueueScheduler) Add(name string, f func()) {
	q.AddTail(name, f)
}

// AddTail 尾插（当前新增任务添加到队列末尾）
func (q *QueueScheduler) AddTail(name string, f func()) {
	if name != "" && f != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.tasks[name]; !exist {
			var task = &QueueTask{name: name, fn: f}
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
			q.tasks[name].fn = f
		}
	}
}

// AddHead 头插（当前新增任务添加到队列首位）
func (q *QueueScheduler) AddHead(name string, f func()) {
	if name != "" && f != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.tasks[name]; !exist {
			var task = &QueueTask{name: name, fn: f}
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
			q.tasks[name].fn = f
		}
	}
}

// AddAfter 后插队(当前新增任务添加到after任务之后)
func (q *QueueScheduler) AddAfter(name string, f func(), after string) {
	if name != "" && f != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.tasks[name]; !exist {
			var task = &QueueTask{name: name, fn: f}
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
			q.tasks[name].fn = f
		}
	}
}

// AddBefore 前插队(当前新增任务添加到before任务之后)
func (q *QueueScheduler) AddBefore(name string, f func(), before string) {
	if name != "" && f != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.tasks[name]; !exist {
			var task = &QueueTask{name: name, fn: f}
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
			q.tasks[name].fn = f
		}
	}
}

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
