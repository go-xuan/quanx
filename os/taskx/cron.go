package taskx

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/robfig/cron/v3"
)

// CronScheduler 定时任务调度器
type CronScheduler struct {
	mutex *sync.Mutex
	cron  *cron.Cron
	tasks map[string]*CronTask
}

// CronTask 定时任务
type CronTask struct {
	name string
	spec string
	*cron.Entry
}

// Corn 定时任务调度器
func Corn() *CronScheduler {
	//初始化一个定时任务调度程序
	return &CronScheduler{
		mutex: new(sync.Mutex),
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
			cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))),
		),
		tasks: make(map[string]*CronTask),
	}
}

// Start 开始执行定时任务
func (s *CronScheduler) Start() {
	s.cron.Start()
}

// Stop 停止执行定时任务
func (s *CronScheduler) Stop() {
	s.cron.Stop()
}

// Add 添加定时任务
func (s *CronScheduler) Add(name string, spec string, task func()) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// 如果已存在同名任务则先移除再新增
	if entry, ok := s.tasks[name]; ok {
		s.cron.Remove(entry.ID)
	}
	// 新增定时任务
	if entryID, err := s.cron.AddFunc(spec, task); err != nil {
		return err
	} else {
		entry := s.cron.Entry(entryID)
		s.tasks[name] = &CronTask{name: name, spec: spec, Entry: &entry}
	}
	return nil
}

// All 获取所有定时任务
func (s *CronScheduler) All() []*CronTask {
	var tasks []*CronTask
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// Get 获取定时任务
func (s *CronScheduler) Get(name string) *CronTask {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if task, ok := s.tasks[name]; ok {
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

// Info 获取定时任务信息
func (s *CronTask) Info() string {
	return fmt.Sprintf("name:%s spec:%s next:%s", s.name, s.spec, s.Next.Format("2006-01-02 15:04:05"))
}
