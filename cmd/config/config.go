package config

import (
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	ServiceName          string        `default:"bybit-trader" envconfig:"SERVICE_NAME"`
	Symbol               string        `default:"BTCUSD" envconfig:"SYMBOL"`
	MinExecutionDuration time.Duration `default:"1s" envconfig:"MIN_EXECUTION_DURATION"`
	MaxExecutionDuration time.Duration `default:"30s" envconfig:"MAX_EXECUTION_DURATION"`
	BybitAPIKey          string        `required:"1" envconfig:"BYBIT_API_KEY"`
	BybitAPISecret       string        `required:"1" envconfig:"BYBIT_API_SECRET"`

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
