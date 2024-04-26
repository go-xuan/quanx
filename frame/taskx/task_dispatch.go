package taskx

import (
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

type Dispatch struct {
	flag    *atomic.Uint32
	mutex   sync.Mutex
	isolate func()
	tasks   []func()
}

func NewDispatch() *Dispatch {
	return &Dispatch{
		flag:  &atomic.Uint32{},
		mutex: sync.Mutex{},
		tasks: make([]func(), 0),
	}
}

func (d *Dispatch) Execute() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.flag.Load() == 1 && d.isolate != nil {
		d.isolate()
	}
	for _, task := range d.tasks {
		task()
	}
	if d.flag.Load() == 2 && d.isolate != nil {
		d.isolate()
	}
	d.flag.Store(0)
	d.isolate = nil
	d.tasks = make([]func(), 0)
}

func (d *Dispatch) Register(task ...func()) {
	if d.flag.Load() == 3 {
		log.Warn("this dispatch has stopped registration")
		return
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.tasks = append(d.tasks, task...)
}

func (d *Dispatch) MustFirst(task func()) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.flag.Store(1)
	if d.isolate == nil {
		d.isolate = task
	} else {
		log.Warn("isolate task has been registered")
	}
}
func (d *Dispatch) MustLast(task func()) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.flag.Store(2)
	if d.isolate == nil {
		d.isolate = task
	} else {
		log.Warn("isolate task has been registered")
	}
}
