package taskx

import (
	"sync"
)

// 任务处理调度器
type Queue struct {
	mutex *sync.Mutex
	head  *QueueTask            // 首个任务
	tail  *QueueTask            // 末尾任务
	tasks map[string]*QueueTask // 所有任务
}

type QueueTask struct {
	name string
	exec func()     // 当前任务
	next *QueueTask // 指向下一个任务
}

func (t *QueueTask) HasNext() bool {
	return t.next != nil
}

func NewQueue() *Queue {
	return &Queue{
		mutex: new(sync.Mutex),
		tasks: make(map[string]*QueueTask),
	}
}

func (q *Queue) Execute() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if curr := q.head; curr != nil {
		curr.exec()
		delete(q.tasks, curr.name)
		for curr.HasNext() {
			curr = curr.next
			curr.exec()
			delete(q.tasks, curr.name)
		}
		q.head = nil
		q.tail = nil
	}
}

// 头插（当前新增任务添加到队列任务最后）
func (q *Queue) Add(name string, fn func()) {
	q.AddTail(name, fn)
}

// 尾插（当前新增任务添加到队列末尾）
func (q *Queue) AddTail(name string, fn func()) {
	if name == "" || fn == nil {
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if _, ok := q.tasks[name]; !ok {
		var newTask = &QueueTask{name: name, exec: fn}
		if tail := q.tail; tail != nil {
			tail.next = newTask
			q.tail = newTask
		} else {
			q.head = newTask
			q.tail = newTask
		}
		q.tasks[name] = newTask
	}
}

// 头插（当前新增任务添加到队列首位）
func (q *Queue) AddHead(name string, fn func()) {
	if name == "" || fn == nil {
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if _, ok := q.tasks[name]; !ok {
		var newTask = &QueueTask{name: name, exec: fn}
		if head := q.head; head != nil {
			newTask.next = head
			q.head = newTask
		} else {
			q.head = newTask
			q.tail = newTask
		}
		q.tasks[name] = newTask
	}
}

// 插队注册任务(当前新增任务添加到after任务之后)
func (q *Queue) AddJump(name string, fn func(), after string) {
	if name == "" || fn == nil {
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if _, ok := q.tasks[name]; !ok {
		var newTask = &QueueTask{name: name, exec: fn}
		var jump *QueueTask
		if jump, ok = q.tasks[after]; ok && jump.next != nil {
			newTask.next = jump.next
			jump.next = newTask
		} else {
			jump = q.tail
			jump.next = newTask
			q.tail = newTask
		}
		q.tasks[name] = newTask
	}
}
