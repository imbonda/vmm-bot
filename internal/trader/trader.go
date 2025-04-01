package trader

import (
	"context"
	"runtime/debug"

	"github.com/go-kit/log"

	"github.com/imbonda/bybit-vmm-bot/cmd/interfaces"
	"github.com/imbonda/bybit-vmm-bot/pkg/models"
	"github.com/imbonda/bybit-vmm-bot/pkg/utils"
)

type Trader struct {
	exchangeClient  interfaces.ExchangeClient
	symbol          string
	spreadMarginMin float64
	spreadMarginMax float64
	tradeQtyMin     float64
	tradeQtyMax     float64
	logger          log.Logger
}

type NewTraderInput struct {
	ExchangeClient  interfaces.ExchangeClient
	Symbol          string
	SpreadMarginMin float64
	SpreadMarginMax float64
	TradeAmountMin  float64
	TradeAmountMax  float64
	Logger          log.Logger
}

type tradeParams struct {
	shouldTrade bool
	price       float64
	qty         float64
}

func NewTrader(ctx context.Context, input *NewTraderInput) (*Trader, error) {
	return &Trader{
		exchangeClient:  input.ExchangeClient,
		symbol:          input.Symbol,
		spreadMarginMin: input.SpreadMarginMin,
		spreadMarginMax: input.SpreadMarginMax,
		tradeQtyMin:     input.TradeAmountMin,
		tradeQtyMax:     input.TradeAmountMax,
		logger:          input.Logger,
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
	return utils.RandInRange(
		spread.Bid+t.spreadMarginMin*spread.Diff,
		spread.Bid+t.spreadMarginMax*spread.Diff,
	)
}

func (t *Trader) getRandQty(ctx context.Context) float64 {
	return utils.RandInRange(t.tradeQtyMin, t.tradeQtyMax)
}
