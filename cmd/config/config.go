package config

import (
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	ServiceName                    string        `default:"bybit-trader" envconfig:"SERVICE_NAME"`
	Symbol                         string        `default:"BTCUSD" envconfig:"SYMBOL"`
	IntervalExecutionDuration      time.Duration `default:"60s" envconfig:"INTERVAL_EXECUTION_DURATION"`
	NumOfTradeIterationsInInterval int           `default:"2" envconfig:"NUM_OF_TRADE_ITERATIONS_IN_INTERVAL"`
	BybitAPIKey                    string        `required:"1" envconfig:"BYBIT_API_KEY"`
	BybitAPISecret                 string        `required:"1" envconfig:"BYBIT_API_SECRET"`

	GraceFullShutdown time.Duration `default:"5s" envconfig:"GRACE_FULL_SHUTDOWN"`

	logger log.Logger
}

func LoadConfig(cfg *Configuration) error {
	return envconfig.Process("", cfg)
}

func (cfg *Configuration) GetLogger() log.Logger {
	if cfg.logger == nil {
		cfg.logger = log.With(log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout)),
			"ts", log.DefaultTimestampUTC, "name", cfg.ServiceName, "symbol", cfg.Symbol)
	}
	return cfg.logger
}
