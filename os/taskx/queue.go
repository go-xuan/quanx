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
	id   string     // 任务ID
	fn   func()     // 当前任务
	prev *QueueTask // 指向上一个任务
	next *QueueTask // 指向下一个任务
}

func (t *QueueTask) HasNext() bool {
	return t.next != nil
}

func (t *QueueTask) Name() string {
	return t.id
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
		delete(q.tasks, q.head.id)
		q.head = q.head.next
	}
	if len(q.tasks) == 0 {
		q.head = nil
		q.tail = nil
	}
}

// Add 新增队列任务（默认尾插）
func (q *QueueScheduler) Add(fn func(), id string) {
	q.AddTail(fn, id)
}

// AddTail 尾插（当前新增任务添加到队列末尾）
func (q *QueueScheduler) AddTail(fn func(), id string) {
	if id != "" && fn != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.tasks[id]; !exist {
			var task = &QueueTask{id: id, fn: fn}
			if tail := q.tail; tail != nil {
				task.prev = tail
				tail.next = task
				q.tail = task
			} else {
				q.head = task
				q.tail = task
			}
			q.tasks[id] = task
		} else {
			q.tasks[id].fn = fn
		}
	}
}

// AddHead 头插（当前新增任务添加到队列首位）
func (q *QueueScheduler) AddHead(fn func(), id string) {
	if id != "" && fn != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.tasks[id]; !exist {
			var task = &QueueTask{id: id, fn: fn}
			if head := q.head; head != nil {
				head.prev = task
				task.next = head
				q.head = task
			} else {
				q.head = task
				q.tail = task
			}
			q.tasks[id] = task
		} else {
			q.tasks[id].fn = fn
		}
	}
}

// AddAfter 后插队(当前新增任务添加到after任务之后)
func (q *QueueScheduler) AddAfter(fn func(), id, afterId string) {
	if id != "" && fn != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.tasks[id]; !exist {
			var task = &QueueTask{id: id, fn: fn}
			if target, ok := q.tasks[afterId]; ok && target.next != nil {
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
			q.tasks[id] = task
		} else {
			q.tasks[id].fn = fn
		}
	}
}

// AddBefore 前插队(当前新增任务添加到before任务之后)
func (q *QueueScheduler) AddBefore(fn func(), id string, beforeId string) {
	if id != "" && fn != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.tasks[id]; !exist {
			var task = &QueueTask{id: id, fn: fn}
			if target, ok := q.tasks[beforeId]; ok && target.prev != nil {
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
			q.tasks[id] = task
		} else {
			q.tasks[id].fn = fn
		}
	}
}

func (q *QueueScheduler) Remove(id string) {
	if id != "" {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if task, ok := q.tasks[id]; ok {
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
			delete(q.tasks, id)
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
