package taskx

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// 定时任务调度器
type CronScheduler struct {
	cron  *cron.Cron
	mutex sync.Mutex
	tasks map[string]*CronTask
}

// 定时任务
type CronTask struct {
	Id   string
	Spec string
	Next time.Time
	*cron.Entry
}

// 定时任务调度器
func NewCronScheduler() *CronScheduler {
	//初始化一个定时任务调度程序
	return &CronScheduler{
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
			cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))),
		),
		mutex: sync.Mutex{},
		tasks: make(map[string]*CronTask),
	}
}

func (s *CronScheduler) Start() {
	s.cron.Start()
}

func (s *CronScheduler) Stop() {
	s.cron.Stop()
}

func (s *CronScheduler) AddTask(id string, spec string, task func()) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// 定时任务如果已存在则先移除再新增
	if entry, ok := s.tasks[id]; ok {
		s.cron.Remove(entry.ID)
	}
	// 新增定时任务
	if entryID, err := s.cron.AddFunc(spec, task); err != nil {
		return err
	} else {
		entry := s.cron.Entry(entryID)
		s.tasks[id] = &CronTask{Id: id, Spec: spec, Entry: &entry}
	}
	return nil
}

func (s *CronScheduler) Tasks() []*CronTask {
	var tasks []*CronTask
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *CronScheduler) GetTask(id string) *CronTask {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if task, ok := s.tasks[id]; ok {
		return task
	}
	return nil
}

func (s *CronScheduler) RemoveTask(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if entry, ok := s.tasks[id]; ok {
		s.cron.Remove(entry.ID)
	}
}

func (s *CronTask) Info() string {
	return fmt.Sprintf("id:%s spec:%s next:%s", s.Id, s.Spec, s.Next.Format("2006-01-02 15:04:05"))
}
