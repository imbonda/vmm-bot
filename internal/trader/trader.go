package trader

import (
	"context"
	"fmt"
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
	spreadMarginLower float64
	spreadMarginUpper float64
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
	SpreadMarginLower float64
	SpreadMarginUpper float64
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
		spreadMarginLower: input.SpreadMarginLower,
		spreadMarginUpper: input.SpreadMarginUpper,
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
	err := t.cancelAllOrders(ctx)
	if err != nil {
		return nil, err
	}
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

func (t *Trader) cancelAllOrders(ctx context.Context) error {
	return t.exchangeClient.CancelAllOrders(ctx, t.symbol)
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
	ticker, err := t.exchangeClient.GetLastTicker(ctx, t.symbol)
	if err != nil {
		return nil, err
	}
	lastPrice, err := ticker.Price()
	if err != nil {
		return nil, err
	}
	spread, err := ticker.Spread()
	if err != nil {
		return nil, err
	}
	oracleTicker, err := t.priceOracleClient.GetLastTicker(ctx, t.oracleSymbol)
	if err != nil {
		return nil, err
	}
	oraclePrice, err := oracleTicker.Price()
	if err != nil {
		return nil, err
	}
	price, err := t.getRandPriceInSpread(ctx, spread, lastPrice, oraclePrice)
	if err != nil {
		return nil, err
	}
	qty := t.getRandQty(ctx)
	return &tradeParams{
		shouldTrade: true,
		price:       utils.FormatFloatToString(price, t.priceDecimals),
		qty:         utils.FormatFloatToString(qty, t.amountDecimals),
	}, nil
}

func (t *Trader) getRandPriceInSpread(_ context.Context, spread *models.Spread, lastPrice float64, oraclePrice float64) (float64, error) {
	// Oracle candle height range
	oracleLowerLimit := oraclePrice * (1 - t.candleHeight/2)
	oracleUpperLimit := oraclePrice * (1 + t.candleHeight/2)

	// Candle height range
	lowerLimit := lastPrice * (1 - t.candleHeight/2)
	upperLimit := lastPrice * (1 + t.candleHeight/2)

	var margin *models.Spread
	if spread.Diff() < 0 {
		clone := spread.Clone()
		clone.Ask = clone.Bid * (1 + t.candleHeight)
		margin = clone.MarginSpread(t.spreadMarginLower, t.spreadMarginUpper)
	} else {
		margin = spread.MarginSpread(t.spreadMarginLower, t.spreadMarginUpper)
	}

	var min, max float64

	switch {
	case margin.Contains(oracleLowerLimit, oracleUpperLimit):
		min, max = oracleLowerLimit, oracleUpperLimit
	case margin.Contains(oracleLowerLimit):
		min, max = oracleLowerLimit, margin.Ask
	case margin.Contains(oracleUpperLimit):
		min, max = margin.Bid, oracleUpperLimit

	case margin.Contains(lowerLimit, upperLimit):
		min, max = lowerLimit, upperLimit
	case margin.Contains(lowerLimit):
		min, max = lowerLimit, margin.Ask
	case margin.Contains(upperLimit):
		min, max = margin.Bid, upperLimit

	// In case the margins are too big consider the spread itself ignoring margins.

	case spread.Contains(oracleLowerLimit, oracleUpperLimit):
		min, max = oracleLowerLimit, oracleUpperLimit
	case spread.Contains(oracleLowerLimit):
		min, max = oracleLowerLimit, spread.Ask
	case spread.Contains(oracleUpperLimit):
		min, max = spread.Bid, oracleUpperLimit

	case spread.Contains(lowerLimit, upperLimit):
		min, max = lowerLimit, upperLimit
	case spread.Contains(lowerLimit):
		min, max = lowerLimit, spread.Ask
	case spread.Contains(upperLimit):
		min, max = spread.Bid, upperLimit

	case spread.Above(oraclePrice):
		min, max = spread.Bid, spread.Bid*(1+t.candleHeight)
	case spread.Below(oraclePrice):
		min, max = spread.Ask*(1-t.candleHeight), spread.Ask

	default:
		min, max = spread.Bid, spread.Ask
	}

	if min > max {
		return 0, fmt.Errorf(
			"unexpected price range. min: %f, max: %f, oraclePrice: %f, price: %f, ask: %f, bid: %f",
			min,
			max,
			oraclePrice,
			lastPrice,
			spread.Ask,
			spread.Bid,
		)
	}

	return utils.RandInRange(min, max), nil
}

func (t *Trader) getRandQty(_ context.Context) float64 {
	return utils.RandInRange(t.tradeQtyMin, t.tradeQtyMax)
}
