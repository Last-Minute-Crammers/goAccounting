package initialize

import (
	"time"

	"github.com/go-co-op/gocron"
)

type _scheduler struct {
}

func (m *_scheduler) initScheduler() error {
	Scheduler = gocron.NewScheduler(time.Local)
	Scheduler.StartAsync()
	return nil
}
