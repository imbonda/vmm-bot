package models

import (
	"time"

	"github.com/go-kit/log"

	"github.com/imbonda/bybit-vmm-bot/cmd/interfaces"
)

type TradeConfig struct {
	Symbol          string
	SpreadMarginMin float64
	SpreadMarginMax float64
	TradeAmountMin  float64
	TradeAmountMax  float64
}

type ExecutorConfig struct {
	IntervalExecutionDuration      time.Duration
	NumOfTradeIterationsInInterval int
	ListenAddress                  string
}

type NewTraderServiceInput struct {
	ExchangeClient interfaces.ExchangeClient
	Trade          TradeConfig
	Executor       ExecutorConfig
	Logger         log.Logger
}
