package taskx

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var CornParser = cron.NewParser(
	cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)

// JobWrapper 任务包装器
type JobWrapper func(name, spec string, job func()) func()

// CronScheduler 定时任务调度器
type CronScheduler struct {
	mutex *sync.Mutex         // 互斥锁
	cron  *cron.Cron          // corn对象
	tasks map[string]*CornJob // 定时任务
	wraps []JobWrapper        // 定时任务包装器
}

type CornJob struct {
	cron.Entry
	name string // 定时任务名
	spec string // 任务定时表达式
	do   func() // 任务执行方法
}

func (j *CornJob) Run() {
	start := time.Now()
	fmt.Printf("%s start do!\n", j.name)
	j.do()
	fmt.Printf("%s finish do, duration: %dms\n", j.name, time.Since(start).Milliseconds())
}

// Info 获取定时任务信息
func (j *CornJob) Info() string {
	return fmt.Sprintf("name:%-20s spec:%-20s prev:%s   next:%s",
		j.name, j.spec,
		j.Prev.Format("2006-01-02 15:04:05"),
		j.Next.Format("2006-01-02 15:04:05"),
	)
}

// Corn 定时任务调度器
func Corn(warps ...JobWrapper) *CronScheduler {
	var options = []cron.Option{
		cron.WithParser(CornParser),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
		cron.WithLogger(cron.DefaultLogger),
	}
	//初始化一个定时任务调度程序
	return &CronScheduler{
		mutex: new(sync.Mutex),
		cron:  cron.New(options...),
		tasks: make(map[string]*CornJob),
		wraps: warps,
	}
}

// Add 添加定时任务
func (s *CronScheduler) Add(name string, spec string, job func()) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// 如果已存在同名任务则先移除再新增
	if cronTask, ok := s.tasks[name]; ok {
		s.cron.Remove(cronTask.ID)
	}
	// 新增定时任务
	var cornJob = new(CornJob)
	cornJob.name = name
	cornJob.spec = spec
	if s.wraps != nil {
		for _, wrap := range s.wraps {
			job = wrap(name, spec, job)
		}
	}
	cornJob.do = job

	if entryID, err := s.cron.AddJob(spec, cornJob); err != nil {
		return err
	} else {
		cornJob.ID = entryID
		s.tasks[name] = cornJob
	}
	return nil
}

func (s *CronScheduler) WrapJob(name, spec string, task func()) func() {
	if s.wraps != nil {
		for _, wrap := range s.wraps {
			task = wrap(name, spec, task)
		}
	}
	return task
}

// Start 开始执行定时任务
func (s *CronScheduler) Start() {
	s.cron.Start()
}

// Stop 停止执行定时任务
func (s *CronScheduler) Stop() {
	s.cron.Stop()
}

// All 获取所有定时任务
func (s *CronScheduler) All() []*CornJob {
	var tasks []*CornJob
	for _, task := range s.tasks {
		task.Entry = s.cron.Entry(task.ID)
		tasks = append(tasks, task)
	}
	return tasks
}

// Get 获取定时任务
func (s *CronScheduler) Get(name string) *CornJob {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if task, ok := s.tasks[name]; ok {
		task.Entry = s.cron.Entry(task.ID)
		return task
	}
	return nil
}

// Remove 移除定时任务
func (s *CronScheduler) Remove(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if entry, ok := s.tasks[name]; ok {
		s.cron.Remove(entry.ID)
		delete(s.tasks, name)
	}
}
