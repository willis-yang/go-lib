package scheduler

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	Cron *cron.Cron
}

func MustNewScheduler() *Scheduler {
	return &Scheduler{
		Cron: NewScheduler(),
	}
}

func NewScheduler() *cron.Cron {
	return cron.New(
		cron.WithParser(cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional)),
	)
}

func (s *Scheduler) AddFunc(spec string, cmd func()) {
	_, err := s.Cron.AddFunc(spec, cmd)
	if err != nil {
		panic(fmt.Sprintf("Failed to Add Func Error: %v", err))
	}

}

func (s Scheduler) Start() {
	s.Cron.Start()
}

func (s Scheduler) Stop() {
	s.Cron.Stop()
}
