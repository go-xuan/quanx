package cronx

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

var cronX *cron.Cron

func New() *cron.Cron {
	if cronX == nil {
		//初始化一个定时任务
		cronX = cron.New(cron.WithSeconds(),
			cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
			cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))
	}
	return cronX
}

func AddTask(spec string, execFunc func()) {
	// 给初始化的定时任务指定时间表达式和具体执行的函数
	entryId, err := New().AddFunc(spec, execFunc)
	if err != nil {
		fmt.Println(err)
	}
	//运行定时任务
	go execFunc()
	fmt.Println(time.Now(), entryId, err)
}
