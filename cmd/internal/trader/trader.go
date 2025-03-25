package trader

import (
	"context"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/go-kit/log"

	"github.com/imbonda/bybit-vmm-bot/cmd/interfaces"
)

type Trader struct {
	exchangeClient interfaces.ExchangeClient
	scheduler      gocron.Scheduler
	symbol         string
	logger         log.Logger
}

type NewTraderInput struct {
	Symbol               string
	ExchangeClient       interfaces.ExchangeClient
	MinExecutionDuration time.Duration
	MaxExecutionDuration time.Duration
	Logger               log.Logger
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
		exchangeClient: input.ExchangeClient,
		scheduler:      scheduler,
		symbol:         input.Symbol,
		logger:         input.Logger,
	}
	_, err = scheduler.NewJob(gocron.DurationRandomJob(input.MinExecutionDuration, input.MaxExecutionDuration),
		gocron.NewTask(func(ctx context.Context) {
			return
		}, gocron.WithSingletonMode(gocron.LimitModeReschedule)),
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

func (t *Trader) tradeOnce() error {
	return nil
}
