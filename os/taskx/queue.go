package taskx

import (
	"sync"
)

// QueueScheduler 队列任务处理调度器
type QueueScheduler struct {
	mutex *sync.Mutex
	Head  *QueueTask            // 首个任务
	Tail  *QueueTask            // 末尾任务
	Tasks map[string]*QueueTask // 所有任务
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
		Tasks: make(map[string]*QueueTask),
	}
}

// Execute 执行队列任务
func (q *QueueScheduler) Execute() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for q.Head != nil {
		q.Head.fn()
		delete(q.Tasks, q.Head.name)
		q.Head = q.Head.next
	}
	if len(q.Tasks) == 0 {
		q.Head = nil
		q.Tail = nil
	}
}

// Add 新增队列任务（默认尾插）
func (q *QueueScheduler) Add(name string, fn func()) {
	if name != "" && fn != nil {
		q.AddTail(name, fn)
	}
}

// AddTail 尾插（当前新增任务添加到队列末尾）
func (q *QueueScheduler) AddTail(name string, fn func()) {
	if name != "" && fn != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.Tasks[name]; !exist {
			var add = &QueueTask{name: name, fn: fn}
			if tail := q.Tail; tail != nil {
				add.prev = tail
				tail.next = add
				q.Tail = add
			} else {
				q.Head = add
				q.Tail = add
			}
			q.Tasks[name] = add
		}
	}
}

// AddHead 头插（当前新增任务添加到队列首位）
func (q *QueueScheduler) AddHead(name string, task func()) {
	if name != "" && task != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.Tasks[name]; !exist {
			var add = &QueueTask{name: name, fn: task}
			if head := q.Head; head != nil {
				head.prev = add
				add.next = head
				q.Head = add
			} else {
				q.Head = add
				q.Tail = add
			}
			q.Tasks[name] = add
		}
	}
}

// AddAfter 后插队(当前新增任务添加到after任务之后)
func (q *QueueScheduler) AddAfter(name string, task func(), after string) {
	if name != "" && task != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.Tasks[name]; !exist {
			var add = &QueueTask{name: name, fn: task}
			if target, ok := q.Tasks[after]; ok && target.next != nil {
				add.next = target.next
				add.prev = target
				target.next.prev = add
				target.next = add
			} else {
				target = q.Tail
				target.next = add
				add.prev = target
				q.Tail = add
			}
			q.Tasks[name] = add
		}
	}
}

// AddBefore 前插队(当前新增任务添加到before任务之后)
func (q *QueueScheduler) AddBefore(name string, task func(), before string) {
	if name != "" && task != nil {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if _, exist := q.Tasks[name]; !exist {
			var add = &QueueTask{name: name, fn: task}
			if target, ok := q.Tasks[before]; ok && target.prev != nil {
				add.prev = target.prev
				add.next = target
				target.prev.next = add
				target.prev = add
			} else {
				target = q.Head
				target.prev = add
				add.next = target
				q.Head = add
			}
			q.Tasks[name] = add
		}
	}
}

func (q *QueueScheduler) Remove(name string) {
	if name != "" {
		q.mutex.Lock()
		defer q.mutex.Unlock()
		if task, ok := q.Tasks[name]; ok {
			if task.prev == nil {
				task.next.prev = nil
				q.Head = task.next
			} else if task.next == nil {
				task.prev.next = nil
				q.Tail = task.prev
			} else {
				task.prev.next = task.next
				task.next.prev = task.prev
			}
			delete(q.Tasks, name)
		}
	}
}

// Clear 清除所有队列任务
func (q *QueueScheduler) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.Head = nil
	q.Tail = nil
	q.Tasks = make(map[string]*QueueTask)
}
