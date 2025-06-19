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

func (q *QueueScheduler) Valid() bool {
	return q.head != nil && q.tail != nil && len(q.tasks) > 0
}

// Reset 队列重置
func (q *QueueScheduler) Reset() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.head = nil
	q.tail = nil
	q.tasks = make(map[string]*QueueTask)
}

func (q *QueueScheduler) Exist(name string) bool {
	if _, ok := q.tasks[name]; ok {
		return true
	}
	return false
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

	// 检查队列是否为空
	if q.head == nil {
		log.Info("queue is empty, no tasks to execute")
		return nil
	}

	current := q.head
	for current != nil {
		logger := log.WithField("curr_task", current.name)
		if current.next != nil {
			logger.WithField("next_task", current.next.name)
		}
		// 执行当前任务
		logger.Info("queue task execute")
		if err := current.fn(); err != nil {
			logger.Error("queue task execute error")
			return err
		}

		// 从任务列表中删除当前任务并更新当前任务指针
		delete(q.tasks, current.name)
		current = current.next
	}

	// 队列重置
	q.Reset()
	return nil
}

// Add 新增队列任务（默认尾插）
func (q *QueueScheduler) Add(name string, fn func() error) {
	q.AddTail(name, fn)
}

// 若任务已存在，更新任务函数
func (q *QueueScheduler) existAndUpdate(name string, fn func() error) bool {
	if existTask, exist := q.tasks[name]; exist {
		log.WithField("add_task", name).Errorf("task already exists, updated function")
		existTask.fn = fn
		return true
	}
	return false
}

// AddTail 尾插（当前新增任务添加到队列末尾）
func (q *QueueScheduler) AddTail(name string, fn func() error) {
	logger := log.WithField("add_task", name)
	if name == "" || fn == nil {
		logger.Error("task name or task function is empty")
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.existAndUpdate(name, fn) {
		return
	}

	newTask := &QueueTask{name: name, fn: fn}
	if q.tail != nil {
		// 队列不为空，将新任务添加到尾部
		newTask.prev = q.tail
		q.tail.next = newTask
		q.tail = newTask
	} else {
		// 队列为空，新任务既是头也是尾
		q.head = newTask
		q.tail = newTask
	}

	q.tasks[name] = newTask
	logger.Info("add task to the tail of the queue")
}

// AddHead 头插（当前新增任务添加到队列首位）
func (q *QueueScheduler) AddHead(name string, fn func() error) {
	logger := log.WithField("add_task", name)
	if name == "" || fn == nil {
		logger.Error("task name or task function is empty")
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.existAndUpdate(name, fn) {
		return
	}

	newTask := &QueueTask{name: name, fn: fn}
	if q.head != nil {
		// 队列已有任务，将新任务插入到头部
		q.head.prev = newTask
		newTask.next = q.head
		q.head = newTask
	} else {
		// 队列为空，新任务既是头也是尾
		q.head = newTask
		q.tail = newTask
	}
	q.tasks[name] = newTask
	logger.Info("add task to the head of the queue")
}

// AddAfter 后插队(将新任务添加到目标任务之后)
// 基于XXX后置添加任务
func (q *QueueScheduler) AddAfter(base, name string, fn func() error) {
	logger := log.WithField("add_task", name).WithField("base_task", base)
	if name == "" || fn == nil {
		logger.Error("task name or task function is empty")
		return
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.existAndUpdate(name, fn) {
		return
	}

	newTask := &QueueTask{name: name, fn: fn}
	if curr, ok := q.tasks[base]; ok {
		if curr.next != nil {
			// 目标任务存在且不是队尾，则插入到目标任务之后
			newTask.next = curr.next
			newTask.prev = curr
			curr.next.prev = newTask
			curr.next = newTask
		} else {
			// 插入当前队尾之后
			curr.next = newTask
			newTask.prev = curr
			q.tail = newTask
		}
		q.tasks[name] = newTask
		logger.Info("add task after")
	} else {
		logger.Error("base task does not exist, cannot add task after it")
	}
}

// AddBefore 前插队(将新任务添加到目标任务之后)
func (q *QueueScheduler) AddBefore(base, name string, fn func() error) {
	logger := log.WithField("add_task", name).WithField("base_task", base)
	if name == "" || fn == nil {
		logger.Error("task name or task function is empty")
		return
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.existAndUpdate(name, fn) {
		return
	}

	newTask := &QueueTask{name: name, fn: fn}
	if curr, ok := q.tasks[base]; ok {
		if curr.prev != nil {
			// 目标任务不是队列头部，插入到目标任务之前
			newTask.prev = curr.prev
			newTask.next = curr
			curr.prev.next = newTask
			curr.prev = newTask
		} else {
			// 目标任务是队列头部，新任务成为新的头部
			newTask.next = curr
			curr.prev = newTask
			q.head = newTask
		}
		q.tasks[name] = newTask
		logger.Info("add task before")
	} else {
		logger.Error("base task does not exist, cannot add task before it")
	}
}

// Remove 移除任务
func (q *QueueScheduler) Remove(name string) {
	if name != "" {
		logger := log.WithField("remove_task", name)
		if task, ok := q.tasks[name]; ok {
			q.mutex.Lock()
			defer q.mutex.Unlock()
			if task.prev == nil && task.next == nil {
				// 移除的是队列中唯一的任务
				q.head = nil
				q.tail = nil
			} else if task.prev == nil {
				// 移除头部任务
				q.head = task.next
				if q.head != nil {
					q.head.prev = nil
				}
			} else if task.next == nil {
				// 移除尾部任务
				q.tail = task.prev
				if q.tail != nil {
					q.tail.next = nil
				}
			} else {
				// 移除中间任务
				task.prev.next = task.next
				task.next.prev = task.prev
			}
			delete(q.tasks, name)
			logger.Info("task has been removed")
		} else {
			logger.Error("task does not exist")
		}
	}
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
