package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"github.com/imbonda/bybit-vmm-bot/cmd/interfaces"
	"github.com/imbonda/bybit-vmm-bot/pkg/exchanges"
	"github.com/imbonda/bybit-vmm-bot/pkg/exchanges/biconomy"
	"github.com/imbonda/bybit-vmm-bot/pkg/exchanges/bybit"
	"github.com/imbonda/bybit-vmm-bot/pkg/utils"
)

type Configuration struct {
	Service struct {
		Name          string              `default:"trader" envconfig:"SERVICE_NAME"`
		Orchestration utils.Orchestration `default:"executor" envconfig:"SERVICE_ORCHESTRATION"`
	}

	Executor struct {
		IntervalExecutionDuration      time.Duration `default:"60s" envconfig:"INTERVAL_EXECUTION_DURATION"`
		NumOfTradeIterationsInInterval int           `default:"2" envconfig:"NUM_OF_TRADE_ITERATIONS_IN_INTERVAL"`
		ListenAddress                  string        `default:":8080" envconfig:"LISTEN_ADDRESS"`
	}

	Exchange struct {
		Name  exchanges.Exchange `required:"1" envconfig:"EXCHANGE_NAME"`
		Bybit struct {
			ExchangeAPIKey    string `required:"1" envconfig:"BYBIT_API_KEY"`
			ExchangeAPISecret string `required:"1" envconfig:"BYBIT_API_SECRET"`
		}
		Biconomy struct {
			ExchangeAPIKey    string `required:"1" envconfig:"BICONOMY_API_KEY"`
			ExchangeAPISecret string `required:"1" envconfig:"BICONOMY_API_SECRET"`
		}
	}

	Trade struct {
		Symbol          string  `required:"1" envconfig:"SYMBOL"`
		SpreadMarginMin float64 `default:"0" envconfig:"SPREAD_MARGIN_MIN"`
		SpreadMarginMax float64 `default:"1" envconfig:"SPREAD_MARGIN_MAX"`
		TradeAmountMin  float64 `required:"1" envconfig:"TRADE_AMOUNT_MIN"`
		TradeAmountMax  float64 `required:"1" envconfig:"TRADE_AMOUNT_MAX"`
		PriceDecimals   int     `default:"2" envconfig:"PRICE_DECIMALS_PRECISION"`
		AmountDecimals  int     `default:"2" envconfig:"AMOUNT_DECIMALS_PRECISION"`
	}

	GraceFullShutdown time.Duration `default:"5s" envconfig:"GRACE_FULL_SHUTDOWN"`

	logger log.Logger
}

func LoadConfig(cfg *Configuration) error {
	return envconfig.Process("", cfg)
}

func (cfg *Configuration) GetLogger() log.Logger {
	if cfg.logger == nil {
		cfg.logger = log.With(
			log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout)),
			"ts", log.DefaultTimestampUTC,
			"name", cfg.Service.Name,
			"symbol", cfg.Trade.Symbol,
		)
	}
	return cfg.logger
}

func (cfg *Configuration) GetExchangeClient(ctx context.Context) (interfaces.ExchangeClient, error) {
	logger := cfg.GetLogger()
	switch cfg.Exchange.Name {
	case exchanges.Biconomy:
		apiClient, err := biconomy.NewClient(ctx, &biconomy.NewClientInput{
			APIKey:    cfg.Exchange.Biconomy.ExchangeAPIKey,
			APISecret: cfg.Exchange.Biconomy.ExchangeAPISecret,
			Logger:    logger,
		})
		if err != nil {
			level.Error(logger).Log("msg", "failed to create biconomy client", "err", err)
			return nil, err
		}
		return apiClient, nil
	case exchanges.Bybit:
		apiClient, err := bybit.NewClient(ctx, &bybit.NewClientInput{
			APIKey:    cfg.Exchange.Bybit.ExchangeAPIKey,
			APISecret: cfg.Exchange.Bybit.ExchangeAPISecret,
			Logger:    logger,
		})
		if err != nil {
			level.Error(logger).Log("msg", "failed to create bybit client", "err", err)
			return nil, err
		}
		return apiClient, nil
	default:
		return nil, fmt.Errorf("failed to resolve exchange client: %s", cfg.Exchange)
	}
}
