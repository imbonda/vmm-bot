package utils

import (
	"context"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type iterable interface {
	DoIteration(ctx context.Context) error
}

type IterationsExecutor[I iterable] struct {
	scheduler    *Scheduler
	callee       I
	logger       log.Logger
	lastRunEpoch atomic.Uint64
}

type NewIterationsExecutorInput[I iterable] struct {
	Callee                         I
	IntervalExecutionDuration      time.Duration
	NumOfTradeIterationsInInterval int
	Logger                         log.Logger
}

// the object need to get a generic type that has "DoIteration".
func NewIterationsExecutor[I iterable](ctx context.Context, input *NewIterationsExecutorInput[I]) (*IterationsExecutor[I], error) {
	scheduler := NewScheduler(&NewSchedulerInput{
		IntervalDuration:   input.IntervalExecutionDuration,
		NumTasksInInterval: input.NumOfTradeIterationsInInterval,
		Logger:             input.Logger,
	})
	executor := &IterationsExecutor[I]{
		scheduler:    scheduler,
		callee:       input.Callee,
		logger:       input.Logger,
		lastRunEpoch: atomic.Uint64{},
	}
	scheduler.SetTask(
		func(ctx context.Context) {
			executor.doIteration(ctx)
		},
	)
	return executor, nil
}

func (ie *IterationsExecutor[I]) Start(ctx context.Context) error {
	ie.scheduler.Run(ctx)
	return nil
}

func (ie *IterationsExecutor[I]) Shutdown(ctx context.Context) error {
	ie.scheduler.Stop(ctx)
	return nil
}

func (ie *IterationsExecutor[I]) doIteration(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			ie.logger.Log("msg", "panic recovered in doIteration", "err", r)
			debug.PrintStack()
		}
		ie.lastRunEpoch.Store(uint64(time.Now().Unix()))
	}()

	level.Debug(ie.logger).Log("msg", "starting trade iteration")
	if err := ie.callee.DoIteration(ctx); err != nil {
		ie.logger.Log("msg", "failed to iterate", "err", err)
	} else {
		level.Debug(ie.logger).Log("msg", "trade iteration is done")
	}
}
