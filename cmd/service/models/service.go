package models

import (
	"time"

	"github.com/go-kit/log"

	"github.com/imbonda/vmm-bot/cmd/interfaces"
)

type TradeConfig struct {
	Symbol          string
	OracleSymbol    string
	CandleHeight    float64
	SpreadMarginMin float64
	SpreadMarginMax float64
	TradeAmountMin  float64
	TradeAmountMax  float64
	PriceDecimals   int
	AmountDecimals  int
}

type ExecutorConfig struct {
	IntervalExecutionDuration      time.Duration
	NumOfTradeIterationsInInterval int
	ListenAddress                  string
}

type NewTraderServiceInput struct {
	ExchangeClient    interfaces.ExchangeClient
	PriceOracleClient interfaces.ExchangeClient
	Trade             TradeConfig
	Executor          ExecutorConfig
	Logger            log.Logger
}
