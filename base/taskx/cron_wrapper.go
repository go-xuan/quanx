package taskx

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/ipx"
	"github.com/go-xuan/quanx/extra/redisx"
)

// CronJobWrapper 定时任务包装器
type CronJobWrapper func(name, spec string, job func(context.Context)) func(context.Context)

// LockWarp 分布式锁装饰器，定时任务执行时必须争抢到锁才执行
func LockWarp(name, spec string, job func(context.Context)) func(context.Context) {
	return func(ctx context.Context) {
		var logger = log.WithField("job_name", name).WithField("job_spec", spec)
		var key = "cron_job_lock:" + name
		var expiration = ParseDurationBySpec(spec) - time.Millisecond
		var host = ipx.GetLocalIP()
		if ok, err := redisx.GetInstance().SetNX(ctx, key, host, expiration).Result(); ok && err == nil { // 获取到锁才执行
			logger.WithField("host", host).Info("ready to execute")
			job(ctx) // 执行
		} else {
			host, _ = redisx.GetInstance().Get(ctx, key).Result()
			logger.WithField("host", host).Error("already executed by other scheduler")
		}
	}
}

// DurationWarp 执行耗时装饰器，计算定时任务执行耗费时间，单位：毫秒/ms
func DurationWarp(name, spec string, job func(context.Context)) func(context.Context) {
	return func(ctx context.Context) {
		var logger = log.WithField("job_name", name).WithField("job_spec", spec)
		var start = time.Now()
		logger.Info("cron job execute start")
		job(ctx) // 执行
		logger.WithField("duration", time.Since(start).String()).Info("cron job execute finish")
	}
}
