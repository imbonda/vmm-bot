package executor

import (
	"context"

	"github.com/go-kit/log"

	"github.com/imbonda/bybit-vmm-bot/cmd/interfaces"
	"github.com/imbonda/bybit-vmm-bot/cmd/service/models"
	"github.com/imbonda/bybit-vmm-bot/internal/trader"
	"github.com/imbonda/bybit-vmm-bot/pkg/utils"
)

type traderExecutor struct {
	traderClient     interfaces.Trader
	intervalExecutor *utils.IterationsExecutor[*trader.Trader]
	logger           log.Logger
}

func NewTraderService(ctx context.Context, input *models.NewTraderServiceInput) (interfaces.TraderService, error) {
	traderClient, err := trader.NewTrader(ctx, &trader.NewTraderInput{
		ExchangeClient:    input.ExchangeClient,
		PriceOracleClient: input.PriceOracleClient,
		Symbol:            input.Trade.Symbol,
		OracleSymbol:      input.Trade.OracleSymbol,
		CandleHeight:      input.Trade.CandleHeight,
		SpreadMarginMin:   input.Trade.SpreadMarginMin,
		SpreadMarginMax:   input.Trade.SpreadMarginMax,
		TradeAmountMin:    input.Trade.TradeAmountMin,
		TradeAmountMax:    input.Trade.TradeAmountMax,
		PriceDecimals:     input.Trade.PriceDecimals,
		AmountDecimals:    input.Trade.AmountDecimals,
		Logger:            input.Logger,
	})
	if err != nil {
		return nil, err
	}
	executor, err := utils.NewIterationsExecutor(
		ctx,
		&utils.NewIterationsExecutorInput[*trader.Trader]{
			Callee:                         traderClient,
			IntervalExecutionDuration:      input.Executor.IntervalExecutionDuration,
			NumOfTradeIterationsInInterval: input.Executor.NumOfTradeIterationsInInterval,
			Logger:                         input.Logger,
		})
	if err != nil {
		return nil, err
	}
	return &traderExecutor{
		traderClient:     traderClient,
		intervalExecutor: executor,
		logger:           input.Logger,
	}, nil
}

func (s *traderExecutor) Start(ctx context.Context) error {
	return s.intervalExecutor.Start(ctx)
}

func (s *traderExecutor) Shutdown(ctx context.Context) error {
	return s.intervalExecutor.Shutdown(ctx)
}
