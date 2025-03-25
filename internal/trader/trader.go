package trader

import (
	"context"
	"math/rand"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/imbonda/bybit-vmm-bot/cmd/interfaces"
	"github.com/imbonda/bybit-vmm-bot/pkg/models"
)

type Trader struct {
	exchangeClient                 interfaces.ExchangeClient
	scheduler                      gocron.Scheduler
	symbol                         string
	numOfTradeIterationsInInterval int
	intervalExecutionDuration      time.Duration
	logger                         log.Logger
	averageDuration                atomic.Int64
	runsCounter                    atomic.Uint64
	lastRunEpoch                   atomic.Uint64
}

type NewTraderInput struct {
	Symbol                         string
	ExchangeClient                 interfaces.ExchangeClient
	IntervalExecutionDuration      time.Duration
	NumOfTradeIterationsInInterval int
	Logger                         log.Logger
}

type spreadRange struct {
	ask    float64
	bid    float64
	spread float64
}

type shouldTradeOutput struct {
	shouldTrade bool
	spread      float64
}

func NewTrader(ctx context.Context, input *NewTraderInput) (*Trader, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	trader := &Trader{
		exchangeClient:                 input.ExchangeClient,
		scheduler:                      scheduler,
		symbol:                         input.Symbol,
		numOfTradeIterationsInInterval: input.NumOfTradeIterationsInInterval,
		intervalExecutionDuration:      input.IntervalExecutionDuration,
		logger:                         input.Logger,
		averageDuration:                atomic.Int64{},
		runsCounter:                    atomic.Uint64{},
		lastRunEpoch:                   atomic.Uint64{},
	}
	_, err = scheduler.NewJob(gocron.DurationJob(input.IntervalExecutionDuration),
		gocron.NewTask(trader.doInterval), gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		return nil, err
	}
	return trader, nil
}

func (t *Trader) Start(ctx context.Context) error {
	t.scheduler.Start()
	return nil
}

func (t *Trader) Shutdown(ctx context.Context) error {
	return t.scheduler.Shutdown()
}

func (t *Trader) doInterval(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			t.logger.Log("msg", "panic recovered", "err", r)
			debug.PrintStack()
		}
		t.lastRunEpoch.Store(uint64(time.Now().Unix()))
	}()
	t.runsCounter.Add(1)
	runCounter := t.runsCounter.Load()
	level.Debug(t.logger).Log("msg", "starting trade interval", "interval", runCounter)
	num := t.numOfTradeIterationsInInterval
	if num <= 0 {
		return
	}
	totalDuration := t.intervalExecutionDuration
	startTime := time.Now()

	for i := 0; i < num; i++ {
		// Estimate slack before each run
		lastAvg := t.averageDuration.Load()
		if lastAvg > 0 {
			avgDur := time.Duration(lastAvg)
			elapsed := time.Since(startTime)
			remaining := totalDuration - elapsed
			remainingIterations := num - i

			// Max slack = time left minus estimated time needed for rest of iterations
			maxSlack := remaining - (avgDur * time.Duration(remainingIterations))
			if maxSlack > 0 {
				sleepTime := time.Duration(rand.Int63n(int64(maxSlack)))
				level.Debug(t.logger).Log("msg", "got random sleep time", "interval", runCounter, "iteration", i+1, "sleepTime", sleepTime.Seconds())
				time.Sleep(sleepTime)
			}
		}

		iterStart := time.Now()
		level.Debug(t.logger).Log("msg", "starting trade iteration", "interval", runCounter, "iteration", i+1)
		if err := t.tradeOnce(ctx); err != nil {
			t.logger.Log("msg", "trade failed", "err", err, "interval", runCounter, "iteration", i+1)
		} else {
			level.Debug(t.logger).Log("msg", "trade iteration is done", "interval", runCounter, "iteration", i+1)
		}
		iterDur := time.Since(iterStart)

		// Update moving average
		prevAvg := time.Duration(t.averageDuration.Load())
		if prevAvg == 0 {
			t.averageDuration.Store(iterDur.Nanoseconds())
		} else {
			newAvg := (prevAvg*9 + iterDur) / 10
			t.averageDuration.Store(int64(newAvg))
		}
	}
}

func (t *Trader) shouldTrade(ctx context.Context) (*shouldTradeOutput, error) {
	orderBook, err := t.exchangeClient.GetOrderBook(ctx, t.symbol)
	if err != nil {
		return nil, err
	}
	spread, err := orderBook.Spread()
	if err != nil {
		return nil, err
	}
	_ = spread
	return nil, nil
}

func (t *Trader) tradeOnce(ctx context.Context) error {
	spreadRange, err := t.getSpreadRange(ctx)
	if err != nil {
		return err
	}
	price, err := t.getRandPriceInSpread(ctx, spreadRange)
	if err != nil {
		return err
	}
	qty, err := t.getRandQty(ctx)
	if err != nil {
		return err
	}
	err = t.exchangeClient.PlaceOrder(ctx, &models.Order{
		Symbol: t.symbol,
		Action: models.Buy,
		Price:  price,
		Qty:    qty,
	})
	if err != nil {
		return err
	}
	return nil

}

func (t *Trader) getSpreadRange(ctx context.Context) (*spreadRange, error) {
	orderBook, err := t.exchangeClient.GetOrderBook(ctx, t.symbol)
	if err != nil {
		return nil, err
	}
	ask, err := orderBook.Ask()
	if err != nil {
		return nil, err
	}
	bid, err := orderBook.Bid()
	if err != nil {
		return nil, err
	}
	spread, err := orderBook.Spread()
	if err != nil {
		return nil, err
	}
	return &spreadRange{
		ask:    ask,
		bid:    bid,
		spread: spread,
	}, nil
}

func (t *Trader) getRandPriceInSpread(ctx context.Context, spreadRange *spreadRange) (float64, error) {
	price := spreadRange.bid + (spreadRange.ask-spreadRange.bid)/2
	return price, nil
}

func (t *Trader) getRandQty(ctx context.Context) (float64, error) {
	return 0, nil
}
