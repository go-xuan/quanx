package taskx

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/robfig/cron/v3"
)

// 定时任务调度器
type CronScheduler struct {
	mutex *sync.Mutex
	cron  *cron.Cron
	tasks map[string]*CronTask
}

// 定时任务
type CronTask struct {
	name string
	spec string
	*cron.Entry
}

// 定时任务调度器
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

// 开始执行定时任务
func (s *CronScheduler) Start() {
	s.cron.Start()
}

// 停止执行定时任务
func (s *CronScheduler) Stop() {
	s.cron.Stop()
}

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

func (s *CronScheduler) All() []*CronTask {
	var tasks []*CronTask
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *CronScheduler) Get(name string) *CronTask {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if task, ok := s.tasks[name]; ok {
		return task
	}
	return nil
}

func (s *CronScheduler) Remove(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if entry, ok := s.tasks[name]; ok {
		s.cron.Remove(entry.ID)
	}
}

func (s *CronTask) Info() string {
	return fmt.Sprintf("name:%s spec:%s next:%s", s.name, s.spec, s.Next.Format("2006-01-02 15:04:05"))
}
