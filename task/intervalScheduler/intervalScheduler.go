package intervalScheduler

import (
	"errors"
	"sync/atomic"
	"time"
)

type Scheduler struct {
	Interval int64
	flag     atomic.Bool
	fn       func()
}

func MustNewScheduler(Interval int64) (*Scheduler, error) {
	if Interval <= 0 {
		return nil, errors.New("interval must be greater than zero")
	}
	return &Scheduler{
		Interval: Interval,
		flag:     atomic.Bool{},
	}, nil
}

func (s *Scheduler) AddFunc(fn func()) {
	s.fn = fn
}

func (s *Scheduler) Start() { // 改为指针接收器
	s.flag.Store(false) // 原子操作设置标志
	for {
		if s.flag.Load() { // 原子操作读取标志
			break
		}
		if s.fn == nil { // 避免 nil 函数调用 panic
			time.Sleep(time.Duration(s.Interval) * time.Second)
			continue
		}

		start := time.Now()
		s.fn()
		elapsed := time.Since(start)
		sleepDuration := time.Duration(s.Interval)*time.Second - elapsed
		if sleepDuration > 0 { // 仅在执行耗时小于间隔时休眠
			time.Sleep(sleepDuration)
		}
	}
}

func (s *Scheduler) Stop() { // 改为指针接收器
	s.flag.Store(true) // 原子操作设置标志
}
