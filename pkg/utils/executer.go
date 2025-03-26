package utils

import (
	"context"
	"math/rand"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type iterable interface {
	DoIteration(ctx context.Context) error
}

type IterationsExecutor[I iterable] struct {
	scheduler                      gocron.Scheduler
	callee                         I
	numOfTradeIterationsInInterval int
	intervalExecutionDuration      time.Duration
	logger                         log.Logger
	averageDuration                atomic.Int64
	runsCounter                    atomic.Uint64
	lastRunEpoch                   atomic.Uint64
}

type NewIterationsExecutorInput[I iterable] struct {
	Callee                         I
	IntervalExecutionDuration      time.Duration
	NumOfTradeIterationsInInterval int
	Logger                         log.Logger
}

// the object need to get a generic type that has "DoIteration".
func NewIterationsExecutor[I iterable](ctx context.Context, input *NewIterationsExecutorInput[I]) (*IterationsExecutor[I], error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	executor := &IterationsExecutor[I]{
		scheduler:                      scheduler,
		callee:                         input.Callee,
		numOfTradeIterationsInInterval: input.NumOfTradeIterationsInInterval,
		intervalExecutionDuration:      input.IntervalExecutionDuration,
		logger:                         input.Logger,
		averageDuration:                atomic.Int64{},
		runsCounter:                    atomic.Uint64{},
		lastRunEpoch:                   atomic.Uint64{},
	}
	_, err = scheduler.NewJob(
		gocron.DurationJob(input.IntervalExecutionDuration),
		gocron.NewTask(executor.doInterval),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		return nil, err
	}
	return executor, nil
}

func (ie *IterationsExecutor[I]) Start(ctx context.Context) error {
	ie.doInterval(ctx)
	ie.scheduler.Start()
	return nil
}

func (ie *IterationsExecutor[I]) Shutdown(ctx context.Context) error {
	return ie.scheduler.Shutdown()
}

func (ie *IterationsExecutor[I]) doInterval(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			ie.logger.Log("msg", "panic recovered in doInterval", "err", r)
			debug.PrintStack()
		}
		ie.lastRunEpoch.Store(uint64(time.Now().Unix()))
	}()
	ie.runsCounter.Add(1)
	runCounter := ie.runsCounter.Load()
	level.Debug(ie.logger).Log("msg", "starting trade interval", "interval", runCounter)
	num := ie.numOfTradeIterationsInInterval
	if num <= 0 {
		return
	}
	totalDuration := ie.intervalExecutionDuration
	startTime := time.Now()

	for i := 0; i < num; i++ {
		// Estimate slack before each run
		lastAvg := ie.averageDuration.Load()
		if lastAvg > 0 {
			avgDur := time.Duration(lastAvg)
			elapsed := time.Since(startTime)
			remaining := totalDuration - elapsed
			remainingIterations := num - i

			// Max slack = time left minus estimated time needed for rest of iterations
			maxSlack := remaining - (avgDur * time.Duration(remainingIterations))
			if maxSlack > 0 {
				sleepTime := time.Duration(rand.Int63n(int64(maxSlack)))
				level.Debug(ie.logger).Log("msg", "got random sleep time", "interval", runCounter, "iteration", i+1, "sleepTime", sleepTime.Seconds())
				time.Sleep(sleepTime)
			}
		}

		iterStart := time.Now()
		level.Debug(ie.logger).Log("msg", "starting trade iteration", "interval", runCounter, "iteration", i+1)
		if err := ie.callee.DoIteration(ctx); err != nil {
			ie.logger.Log("msg", "failed to iterate", "err", err, "interval", runCounter, "iteration", i+1)
		} else {
			level.Debug(ie.logger).Log("msg", "trade iteration is done", "interval", runCounter, "iteration", i+1)
		}
		iterDur := time.Since(iterStart)

		// Update moving average
		prevAvg := time.Duration(ie.averageDuration.Load())
		if prevAvg == 0 {
			ie.averageDuration.Store(iterDur.Nanoseconds())
		} else {
			newAvg := (prevAvg*9 + iterDur) / 10
			ie.averageDuration.Store(int64(newAvg))
		}
	}
}
