package taskx

import (
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

const (
	normal = iota
	isolateBefore
	isolateAfter
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
	if d.flag.Load() == isolateBefore && d.isolate != nil {
		d.isolate()
	}
	for _, task := range d.tasks {
		task()
	}
	if d.flag.Load() == isolateAfter && d.isolate != nil {
		d.isolate()
	}
	d.flag.Store(normal)
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

func (d *Dispatch) Before(task func()) {
	d.setIsolate(task, isolateBefore)
}

func (d *Dispatch) After(task func()) {
	d.setIsolate(task, isolateAfter)
}

func (d *Dispatch) setIsolate(task func(), flag uint32) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.flag.Load() == normal && d.isolate == nil {
		d.isolate = task
		d.flag.Store(flag)
	} else if d.flag.Load() == isolateBefore {
		log.Warn("isolate task has been registered as before task")
	} else if d.flag.Load() == isolateAfter {
		log.Warn("isolate task has been registered as after task")
	}
}
