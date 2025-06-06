package service

import (
	"context"
	"fmt"

	"github.com/go-kit/log/level"

	"github.com/imbonda/vmm-bot/cmd/config"
	"github.com/imbonda/vmm-bot/cmd/interfaces"
	"github.com/imbonda/vmm-bot/cmd/service/executor"
	"github.com/imbonda/vmm-bot/cmd/service/http"
	"github.com/imbonda/vmm-bot/cmd/service/models"
	"github.com/imbonda/vmm-bot/pkg/utils"
)

func GetTraderService(ctx context.Context, cfg *config.Configuration) (interfaces.TraderService, error) {
	logger := cfg.GetLogger()
	exchangeClient, err := cfg.GetExchangeClient(ctx)
	if err != nil {
		level.Error(logger).Log("msg", "failed to create exchange client", "err", err)
		return nil, err
	}
	priceOracleClient, err := cfg.GetPriceOracleClient(ctx)
	if err != nil {
		level.Error(logger).Log("msg", "failed to create price oracle client", "err", err)
		return nil, err
	}
	if cfg.Service.Orchestration == utils.Executor {
		return executor.NewTraderService(ctx, &models.NewTraderServiceInput{
			ExchangeClient:    exchangeClient,
			PriceOracleClient: priceOracleClient,
			Trade:             models.TradeConfig(cfg.Trade),
			Executor:          models.ExecutorConfig(cfg.Executor),
			Logger:            logger,
		})
	} else if cfg.Service.Orchestration == utils.HTTP {
		return http.NewTraderService(ctx, &models.NewTraderServiceInput{
			ExchangeClient:    exchangeClient,
			PriceOracleClient: priceOracleClient,
			Trade:             models.TradeConfig(cfg.Trade),
			Executor:          models.ExecutorConfig(cfg.Executor),
			Logger:            logger,
		})
	} else {
		level.Error(logger).Log("msg", "invalid orchestration", "orchestration", cfg.Service.Orchestration)
		return nil, fmt.Errorf("invalid orchestration")
	}
}
