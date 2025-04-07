package utils

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Scheduler struct {
	intervalDuration   time.Duration
	numTasksInInterval int
	task               func(context.Context)
	logger             log.Logger
	taskChan           chan struct{}
	stopChan           chan struct{}
}

type NewSchedulerInput struct {
	IntervalDuration   time.Duration
	NumTasksInInterval int
	Task               func(context.Context)
	Logger             log.Logger
}

func NewScheduler(input *NewSchedulerInput) *Scheduler {
	return &Scheduler{
		intervalDuration:   input.IntervalDuration,
		numTasksInInterval: input.NumTasksInInterval,
		task:               input.Task,
		logger:             input.Logger,
		taskChan:           make(chan struct{}),
		stopChan:           make(chan struct{}),
	}
}

func (s *Scheduler) SetTask(task func(context.Context)) {
	s.task = task
}

func (s *Scheduler) Stop(ctx context.Context) {
	close(s.stopChan)
}

func (s *Scheduler) Run(ctx context.Context) {
	go s.run(ctx)
}

func (s *Scheduler) run(ctx context.Context) {
	for {
		level.Info(s.logger).Log("msg", "starting run interval")
		start := time.Now()
		// Schedule operations randomly within interval
		for range s.numTasksInInterval {
			go s.scheduleTask()
		}

		// Consume tasks sequentially
		for range s.numTasksInInterval {
			select {
			case <-s.stopChan:
				level.Info(s.logger).Log("msg", "stopped run loop")
				return
			case <-s.taskChan:
				s.task(ctx)
			}
		}

		level.Debug(s.logger).Log("msg", "finished run interval")

		// Measure how long it took to schedule
		elapsed := time.Since(start)
		remaining := s.intervalDuration - elapsed
		if remaining > 0 {
			select {
			case <-s.stopChan:
				level.Info(s.logger).Log("msg", "stopped run loop")
				return
			case <-time.After(remaining):
				continue
			}
		}
	}
}

func (s *Scheduler) scheduleTask() {
	delay := time.Duration(rand.Int63n(int64(s.intervalDuration.Milliseconds()))) * time.Millisecond
	level.Debug(s.logger).Log("msg", "got random sleep time", "sleepTime", delay.Seconds())
	select {
	case <-s.stopChan:
		return
	case <-time.After(delay):
		s.taskChan <- struct{}{}
	}
}
