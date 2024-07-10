package task

import (
	"bwhNotify/logger"
	"bwhNotify/util"
	"sync"

	"github.com/robfig/cron/v3"
)

type Tasker struct {
	c *cron.Cron
	l *logger.Logger
}

var (
	once   sync.Once
	tasker *Tasker
)

func NewTasker(l *logger.Logger) *Tasker {
	once.Do(func() {
		c := cron.New(
			cron.WithSeconds(),
			cron.WithChain(cron.Recover(l)),
			cron.WithChain(cron.SkipIfStillRunning(l)),
			cron.WithLocation(util.CST),
		)
		tasker = &Tasker{
			c: c,
			l: l,
		}
	})
	return tasker
}

func (t *Tasker) Start() error {
	//加载定时任务内容
	_, err := t.c.AddJob("@daily", cron.NewChain().Then(&getBwhStatus{}))
	if err != nil {
		return err
	}
	t.c.Start()
	return nil
}

func (t *Tasker) Stop() {
	t.c.Stop()
}
