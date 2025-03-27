package trader

import (
	"context"
	"runtime/debug"

	"github.com/go-kit/log"

	"github.com/imbonda/bybit-vmm-bot/cmd/interfaces"
	"github.com/imbonda/bybit-vmm-bot/pkg/models"
)

type Trader struct {
	exchangeClient interfaces.ExchangeClient
	symbol         string
	logger         log.Logger
}

type NewTraderInput struct {
	ExchangeClient interfaces.ExchangeClient
	Symbol         string
	Logger         log.Logger
}

type tradeParams struct {
	shouldTrade bool
	price       float64
	qty         float64
}

func NewTrader(ctx context.Context, input *NewTraderInput) (*Trader, error) {
	return &Trader{
		exchangeClient: input.ExchangeClient,
		symbol:         input.Symbol,
		logger:         input.Logger,
	}, nil
}

func (t *Trader) DoIteration(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			t.logger.Log("msg", "panic recovered in iteration", "err", r)
			debug.PrintStack()
		}
	}()
	_, err := t.TradeOnce(ctx)
	return err
}

func (t *Trader) TradeOnce(ctx context.Context) (*models.TradeOnceOutput, error) {
	params, err := t.getTradeParams(ctx)
	if err != nil {
		return nil, err
	}
	err = t.exchangeClient.PlaceOrder(ctx, &models.Order{
		Symbol: t.symbol,
		Action: models.Sell,
		Price:  params.price,
		Qty:    params.qty,
	})
	if err != nil {
		return nil, err
	}
	err = t.exchangeClient.PlaceOrder(ctx, &models.Order{
		Symbol: t.symbol,
		Action: models.Buy,
		Price:  params.price,
		Qty:    params.qty,
	})
	// TODO: what happens if always fails to buy? .. will sell everything
	return &models.TradeOnceOutput{}, err
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
