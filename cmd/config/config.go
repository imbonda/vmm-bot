package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"github.com/imbonda/vmm-bot/cmd/interfaces"
	"github.com/imbonda/vmm-bot/pkg/exchanges"
	"github.com/imbonda/vmm-bot/pkg/exchanges/biconomy"
	"github.com/imbonda/vmm-bot/pkg/exchanges/bybit"
	"github.com/imbonda/vmm-bot/pkg/utils"
)

type Configuration struct {
	Service struct {
		Name              string              `default:"trader" envconfig:"SERVICE_NAME"`
		Orchestration     utils.Orchestration `default:"executor" envconfig:"SERVICE_ORCHESTRATION"`
		GraceFullShutdown time.Duration       `default:"5s" envconfig:"GRACE_FULL_SHUTDOWN"`
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
		Symbol            string  `required:"1" envconfig:"SYMBOL"`
		OracleSymbol      string  `required:"1" envconfig:"ORACLE_SYMBOL"`
		CandleHeight      float64 `required:"1" envconfig:"CANDLE_HEIGHT"`
		SpreadMarginLower float64 `default:"0" envconfig:"SPREAD_MARGIN_LOWER"`
		SpreadMarginUpper float64 `default:"1" envconfig:"SPREAD_MARGIN_UPPER"`
		TradeAmountMin    float64 `required:"1" envconfig:"TRADE_AMOUNT_MIN"`
		TradeAmountMax    float64 `required:"1" envconfig:"TRADE_AMOUNT_MAX"`
		PriceDecimals     int     `default:"3" envconfig:"PRICE_DECIMALS_PRECISION"`
		AmountDecimals    int     `default:"2" envconfig:"AMOUNT_DECIMALS_PRECISION"`
	}

	Log struct {
		Level  string `default:"all" envconfig:"LOGGER_LEVEL"`
		logger log.Logger
	}
}

func LoadConfig(cfg *Configuration) error {
	return envconfig.Process("", cfg)
}

func (cfg *Configuration) GetLogger() log.Logger {
	if cfg.Log.logger == nil {
		logger := log.With(
			log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout)),
			"ts", log.DefaultTimestampUTC,
			"name", cfg.Service.Name,
			"symbol", cfg.Trade.Symbol,
		)
		var logLevel level.Option
		switch cfg.Log.Level {
		case "debug":
			logLevel = level.AllowDebug()
		case "info":
			logLevel = level.AllowInfo()
		case "warn":
			logLevel = level.AllowWarn()
		case "error":
			logLevel = level.AllowError()
		case "none":
			logLevel = level.AllowNone()
		default:
			logLevel = level.AllowAll()
		}
		cfg.Log.logger = level.NewFilter(logger, logLevel)
	}
	return cfg.Log.logger
}

func (cfg *Configuration) GetExchangeClient(ctx context.Context) (interfaces.ExchangeClient, error) {
	switch cfg.Exchange.Name {
	case exchanges.Biconomy:
		return cfg.getBiconomyClient(ctx)
	case exchanges.Bybit:
		return cfg.getBybitClient(ctx)
	default:
		return nil, fmt.Errorf("failed to resolve exchange client: %s", cfg.Exchange)
	}
}

func (cfg *Configuration) GetPriceOracleClient(ctx context.Context) (interfaces.ExchangeClient, error) {
	return cfg.getBybitClient(ctx)
}

func (cfg *Configuration) getBiconomyClient(ctx context.Context) (interfaces.ExchangeClient, error) {
	logger := cfg.GetLogger()
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
}

func (cfg *Configuration) getBybitClient(ctx context.Context) (interfaces.ExchangeClient, error) {
	logger := cfg.GetLogger()
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
}
