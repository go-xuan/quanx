package cronx

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

type TimeFunc func() error

func TimeTask(spec string, timeFunc TimeFunc) {
	//初始化一个定时任务
	c := cron.New(cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
		cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))
	var i = 1
	// 给初始化的定时任务指定时间表达式和具体执行的函数
	entryId, err := c.AddFunc(spec, func() {
		fmt.Printf("时间 : %d 次数 : %d ", time.Now().UnixMilli(), i)
		err := timeFunc
		if err != nil {
			return
		}
		i++
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(time.Now(), entryId, err)
	//运行定时任务
	c.Start()
	select {}
}
