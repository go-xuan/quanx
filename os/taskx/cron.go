package taskx

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/go-xuan/quanx/os/errorx"
)

var (
	_scheduler *CronScheduler
	CornParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
)

// Corn 定时任务调度器
func Corn(warps ...JobWrapper) *CronScheduler {
	if _scheduler == nil {
		var options = []cron.Option{
			cron.WithParser(CornParser),
			cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
			cron.WithLogger(cron.DefaultLogger),
		}
		//初始化一个定时任务调度程序
		_scheduler = &CronScheduler{
			mutex:   new(sync.Mutex),
			status:  initializationStatus,
			cron:    cron.New(options...),
			entries: make(map[string]*CornEntry),
			wraps:   warps,
		}
	}
	return _scheduler
}

// ParseDurationBySpec 解析表达式，计算当前时间和下次执行时间的时间差
func ParseDurationBySpec(spec string) time.Duration {
	if schedule, err := CornParser.Parse(spec); err == nil {
		var now = time.Now()
		return schedule.Next(now).Sub(now)
	}
	return time.Duration(-1)
}

// JobWrapper 任务包装器
type JobWrapper func(name, spec string, job func()) func()

const (
	initializationStatus = iota
	readinessStatus
	runningStatus
	stopStatus
)

// CronScheduler 定时任务调度器
type CronScheduler struct {
	mutex   *sync.Mutex           // 互斥锁
	status  int                   // 调度器状态（0-初始化；1-待运行；2-运行中；3-停止）
	cron    *cron.Cron            // corn对象
	entries map[string]*CornEntry // 定时任务条目
	wraps   []JobWrapper          // 定时任务包装器
}

// Add 添加定时任务
func (s *CronScheduler) Add(name, spec string, job func()) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// 如果已存在同名任务则先移除再新增
	if cronTask, ok := s.entries[name]; ok {
		s.cron.Remove(cronTask.ID)
	}
	// 新增定时任务
	var entry = new(CornEntry)
	entry.name = name
	entry.spec = spec

	// 遍历装饰器，对任务执行方法进行包装
	if s.wraps != nil {
		for _, wrap := range s.wraps {
			job = wrap(name, spec, job)
		}
	}
	entry.do = job
	if entryID, err := s.cron.AddJob(spec, entry); err != nil {
		return err
	} else {
		entry.ID = entryID
		s.entries[name] = entry
	}
	if s.status == initializationStatus {
		s.status = readinessStatus
	}
	return nil
}

// Start 开始执行定时任务
func (s *CronScheduler) Start() error {
	switch s.status {
	case initializationStatus:
		return errorx.New("please add the job first")
	case runningStatus:
		return errorx.New("the cron scheduler already running")
	default:
		s.cron.Start()
		s.status = runningStatus
		return nil
	}
}

// Stop 停止执行定时任务
func (s *CronScheduler) Stop() error {
	switch s.status {
	case initializationStatus, readinessStatus:
		return errorx.New("the cron scheduler is not running yet")
	case stopStatus:
		return errorx.New("the cron scheduler has stopped")
	default:
		s.cron.Stop()
		s.status = stopStatus
		return nil
	}
}

// All 获取所有定时任务
func (s *CronScheduler) All() []*CornEntry {
	var tasks []*CornEntry
	for _, task := range s.entries {
		task.Entry = s.cron.Entry(task.ID)
		tasks = append(tasks, task)
	}
	return tasks
}

// Get 获取定时任务
func (s *CronScheduler) Get(name string) *CornEntry {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if task, ok := s.entries[name]; ok {
		task.Entry = s.cron.Entry(task.ID)
		return task
	}
	return nil
}

// Remove 移除定时任务
func (s *CronScheduler) Remove(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if entry, ok := s.entries[name]; ok {
		s.cron.Remove(entry.ID)
		delete(s.entries, name)
	}
	// 当任务清零则状态值归零
	if len(s.entries) == 0 {
		if s.status == runningStatus {
			s.cron.Stop()
		}
		s.status = initializationStatus
	}
}

type CornEntry struct {
	cron.Entry
	name string // 定时任务名
	spec string // 任务定时表达式
	do   func() // 任务执行方法
}

func (e *CornEntry) Run() {
	e.do()
}

// Info 获取定时任务信息
func (e *CornEntry) Info() string {
	return fmt.Sprintf("name:%-20s spec:%-20s prev:%s   next:%s",
		e.name, e.spec,
		e.Prev.Format("2006-01-02 15:04:05"),
		e.Next.Format("2006-01-02 15:04:05"),
	)
}
