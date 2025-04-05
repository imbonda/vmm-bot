package trader

import (
	"context"
	"math"
	"runtime/debug"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/imbonda/vmm-bot/cmd/interfaces"
	"github.com/imbonda/vmm-bot/pkg/models"
	"github.com/imbonda/vmm-bot/pkg/utils"
)

type Trader struct {
	exchangeClient    interfaces.ExchangeClient
	priceOracleClient interfaces.ExchangeClient
	symbol            string
	oracleSymbol      string
	candleHeight      float64
	spreadMarginMin   float64
	spreadMarginMax   float64
	tradeQtyMin       float64
	tradeQtyMax       float64
	priceDecimals     int
	amountDecimals    int
	logger            log.Logger
}

type NewTraderInput struct {
	ExchangeClient    interfaces.ExchangeClient
	PriceOracleClient interfaces.ExchangeClient
	Symbol            string
	OracleSymbol      string
	CandleHeight      float64
	SpreadMarginMin   float64
	SpreadMarginMax   float64
	TradeAmountMin    float64
	TradeAmountMax    float64
	PriceDecimals     int
	AmountDecimals    int
	Logger            log.Logger
}

type tradeParams struct {
	shouldTrade bool
	price       string
	qty         string
}

func NewTrader(ctx context.Context, input *NewTraderInput) (*Trader, error) {
	return &Trader{
		exchangeClient:    input.ExchangeClient,
		priceOracleClient: input.PriceOracleClient,
		symbol:            input.Symbol,
		oracleSymbol:      input.OracleSymbol,
		candleHeight:      input.CandleHeight,
		spreadMarginMin:   input.SpreadMarginMin,
		spreadMarginMax:   input.SpreadMarginMax,
		tradeQtyMin:       input.TradeAmountMin,
		tradeQtyMax:       input.TradeAmountMax,
		priceDecimals:     input.PriceDecimals,
		amountDecimals:    input.AmountDecimals,
		logger:            input.Logger,
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
	err = t.placeOrder(ctx, params, models.Sell)
	if err != nil {
		return nil, err
	}
	err = t.placeOrder(ctx, params, models.Buy)
	return &models.TradeOnceOutput{}, err
}

func (t *Trader) placeOrder(ctx context.Context, orderInput *tradeParams, action models.OrderAction) error {
	err := t.exchangeClient.PlaceOrder(ctx, &models.Order{
		Symbol: t.symbol,
		Price:  orderInput.price,
		Qty:    orderInput.qty,
		Action: action,
	})
	if err != nil {
		level.Warn(t.logger).Log(
			"msg", "failed order",
			"symbol", t.symbol,
			"action", action,
			"price", orderInput.price,
			"qty", orderInput.qty,
		)
		return err
	}
	level.Info(t.logger).Log(
		"msg", "successful order",
		"symbol", t.symbol,
		"action", action,
		"price", orderInput.price,
		"qty", orderInput.qty,
	)
	return nil
}

func (t *Trader) getTradeParams(ctx context.Context) (*tradeParams, error) {
	ticker, err := t.exchangeClient.GetLatestTicker(ctx, t.symbol)
	if err != nil {
		return nil, err
	}
	spread, err := ticker.Spread()
	if err != nil {
		return nil, err
	}
	oracleTicker, err := t.priceOracleClient.GetLatestTicker(ctx, t.oracleSymbol)
	if err != nil {
		return nil, err
	}
	lastPrice, err := oracleTicker.Price()
	if err != nil {
		return nil, err
	}
	price := t.getRandPriceInSpread(ctx, spread, lastPrice)
	qty := t.getRandQty(ctx)
	return &tradeParams{
		shouldTrade: true,
		price:       utils.FormatFloatToString(price, t.priceDecimals),
		qty:         utils.FormatFloatToString(qty, t.amountDecimals),
	}, nil
}

func (t *Trader) getRandPriceInSpread(_ context.Context, spread *models.Spread, lastPrice float64) float64 {
	// Calculate the intersection range
	lowerLimit := lastPrice * (1 - t.candleHeight)
	upperLimit := lastPrice * (1 + t.candleHeight)

	// Spread bounds
	spreadMin := spread.Bid + t.spreadMarginMin*spread.Diff
	spreadMax := spread.Bid + t.spreadMarginMax*spread.Diff

	// Final intersection range
	min := math.Max(spreadMin, lowerLimit)
	max := math.Min(spreadMax, upperLimit)

	if min >= max {
		min, max = spreadMin, spreadMax
	}
	return utils.RandInRange(min, max)
}

func (t *Trader) getRandQty(_ context.Context) float64 {
	return utils.RandInRange(t.tradeQtyMin, t.tradeQtyMax)
}
