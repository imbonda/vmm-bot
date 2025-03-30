package config

import (
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/kelseyhightower/envconfig"

	"github.com/imbonda/bybit-vmm-bot/pkg/utils"
)

type Configuration struct {
	ServiceName                    string              `default:"trader" envconfig:"SERVICE_NAME"`
	ServiceOrchestration           utils.Orchestration `default:"executor" envconfig:"SERVICE_ORCHESTRATION"`
	IntervalExecutionDuration      time.Duration       `default:"60s" envconfig:"INTERVAL_EXECUTION_DURATION"`
	NumOfTradeIterationsInInterval int                 `default:"2" envconfig:"NUM_OF_TRADE_ITERATIONS_IN_INTERVAL"`
	ListenAddress                  string              `default:":8080" envconfig:"LISTEN_ADDRESS"`
	ExchangeAPIKey                 string              `required:"1" envconfig:"EXCHANGE_API_KEY"`
	ExchangeAPISecret              string              `required:"1" envconfig:"EXCHANGE_API_SECRET"`
	Symbol                         string              `required:"1" envconfig:"SYMBOL"`

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
