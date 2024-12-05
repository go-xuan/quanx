package taskx

import (
	"github.com/go-xuan/quanx/net/httpx"
	log "github.com/sirupsen/logrus"
)

// CronJob 定时任务注册通用接口
type CronJob interface {
	Register()
}

// FuncCronJob 自定义方法类定时任务，实现 CronJob 接口
type FuncCronJob struct {
	Name string // 任务名
	Spec string // 定时表达式
	Do   func() // 执行函数
}

// Register 任务注册
func (t *FuncCronJob) Register() {
	if err := Corn().Add(t.Name, t.Spec, t.Do); err != nil {
		log.WithField("job_name", t.Name).WithField("job_spec", t.Spec).
			Error("register func-cron-job failed: " + err.Error())
	}
}

// HttpCronJob http请求类定时任务，实现 CronJob 接口
type HttpCronJob struct {
	Name     string               // 任务名
	Spec     string               // 定时表达式
	Strategy httpx.ClientStrategy // http客户端类型
	Request  *httpx.Request       // http请求
}

// Register 任务注册
func (t *HttpCronJob) Register() {
	if err := Corn().Add(t.Name, t.Spec, func() {
		if _, err := t.Request.Do(t.Strategy); err != nil {
			log.Error("request failed: ", err)
		}
	}); err != nil {
		log.WithField("job_name", t.Name).WithField("job_spec", t.Spec).
			Error("register http-cron-job failed: ", err)
	}
}
