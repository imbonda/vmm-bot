package trader

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/go-kit/log"

	"github.com/imbonda/bybit-vmm-bot/cmd/interfaces"
	"github.com/imbonda/bybit-vmm-bot/pkg/models"
	"github.com/imbonda/bybit-vmm-bot/pkg/utils"
)

type Trader struct {
	exchangeClient interfaces.ExchangeClient
	executor       *utils.IterationsExecutor[*Trader]
	symbol         string
	logger         log.Logger
}

type NewTraderInput struct {
	ExchangeClient                 interfaces.ExchangeClient
	IntervalExecutionDuration      time.Duration
	NumOfTradeIterationsInInterval int
	Symbol                         string
	Logger                         log.Logger
}

type tradeParams struct {
	shouldTrade bool
	price       float64
	qty         float64
}

func NewTrader(ctx context.Context, input *NewTraderInput) (*Trader, error) {
	trader := &Trader{
		exchangeClient: input.ExchangeClient,
		symbol:         input.Symbol,
		logger:         input.Logger,
	}
	executor, err := utils.NewIterationsExecutor(
		ctx,
		&utils.NewIterationsExecutorInput[*Trader]{
			Callee:                         trader,
			IntervalExecutionDuration:      input.IntervalExecutionDuration,
			NumOfTradeIterationsInInterval: input.NumOfTradeIterationsInInterval,
			Logger:                         input.Logger,
		})
	if err != nil {
		return nil, err
	}
	trader.executor = executor
	return trader, nil
}

func (t *Trader) Start(ctx context.Context) error {
	return t.executor.Start(ctx)
}

func (t *Trader) Shutdown(ctx context.Context) error {
	return t.executor.Shutdown(ctx)
}

func (t *Trader) DoIteration(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			t.logger.Log("msg", "panic recovered in iteration", "err", r)
			debug.PrintStack()
		}
	}()
	return t.tradeOnce(ctx)
}

func (t *Trader) tradeOnce(ctx context.Context) error {
	params, err := t.getTradeParams(ctx)
	if err != nil {
		return err
	}
	err = t.exchangeClient.PlaceOrder(ctx, &models.Order{
		Symbol: t.symbol,
		Action: models.Sell,
		Price:  params.price,
		Qty:    params.qty,
	})
	if err != nil {
		return err
	}
	err = t.exchangeClient.PlaceOrder(ctx, &models.Order{
		Symbol: t.symbol,
		Action: models.Buy,
		Price:  params.price,
		Qty:    params.qty,
	})
	// TODO: what happens if always fails to buy? .. will sell everything
	return err
}

func (t *Trader) getTradeParams(ctx context.Context) (*tradeParams, error) {
	orderBook, err := t.exchangeClient.GetOrderBook(ctx, t.symbol)
	if err != nil {
		return nil, err
	}
	spread, err := orderBook.Spread()
	if err != nil {
		return nil, err
	}
	price := t.getRandPriceInSpread(ctx, spread)
	qty := t.getRandQty(ctx)
	return &tradeParams{
		shouldTrade: true,
		price:       price,
		qty:         qty,
	}, nil
}

func (t *Trader) getRandPriceInSpread(ctx context.Context, spread *models.Spread) float64 {
	price := spread.Bid + (spread.Ask-spread.Bid)/2 /// 100 ...125.. 150
	return price
}

func (t *Trader) getRandQty(ctx context.Context) float64 {
	return 0
}
