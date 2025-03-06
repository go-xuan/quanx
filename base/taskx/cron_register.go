package taskx

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/httpx"
)

// CronJob 定时任务注册通用接口
type CronJob interface {
	Register()
}

// FunctionCronJob 自定义方法类定时任务，实现 CronJob 接口
type FunctionCronJob struct {
	Name string                // 任务名
	Spec string                // 定时表达式
	Do   func(context.Context) // 执行函数
}

// Register 任务注册
func (t *FunctionCronJob) Register() {
	if err := Corn().Add(t.Name, t.Spec, t.Do); err != nil {
		log.WithField("job_name", t.Name).WithField("job_spec", t.Spec).
			Error("register function cron job failed: " + err.Error())
	}
}

// RequestCronJob http请求类定时任务，实现 CronJob 接口
type RequestCronJob struct {
	Name       string            // 任务名
	Spec       string            // 定时表达式
	Request    *httpx.Request    // http请求
	ClientUses []httpx.ClientUse // httpx.Client使用
}

// Register 任务注册
func (t *RequestCronJob) Register() {
	if err := Corn().Add(t.Name, t.Spec, func(context.Context) {
		if _, err := t.Request.Do(t.ClientUses...); err != nil {
			log.Error("do request failed: ", err)
		}
	}); err != nil {
		log.WithField("job_name", t.Name).WithField("job_spec", t.Spec).
			Error("register request cron job failed: ", err)
	}
}
