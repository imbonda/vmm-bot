package models

import (
	"time"

	"github.com/go-kit/log"

	"github.com/imbonda/bybit-vmm-bot/cmd/interfaces"
)

type NewTraderServiceInput struct {
	ExchangeClient interfaces.ExchangeClient
	Symbol         string

	IntervalExecutionDuration      time.Duration
	NumOfTradeIterationsInInterval int
	ListenAddress                  string

	Logger log.Logger
}
