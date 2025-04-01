package service

import (
	"context"
	"fmt"

	"github.com/go-kit/log/level"

	"github.com/imbonda/bybit-vmm-bot/cmd/config"
	"github.com/imbonda/bybit-vmm-bot/cmd/interfaces"
	"github.com/imbonda/bybit-vmm-bot/cmd/service/executor"
	"github.com/imbonda/bybit-vmm-bot/cmd/service/http"
	"github.com/imbonda/bybit-vmm-bot/cmd/service/models"
	"github.com/imbonda/bybit-vmm-bot/pkg/utils"
)

func GetTraderService(ctx context.Context, cfg *config.Configuration) (interfaces.TraderService, error) {
	logger := cfg.GetLogger()
	exchangeClient, err := cfg.GetExchangeClient(ctx)
	if err != nil {
		level.Error(logger).Log("msg", "failed to create bybit client", "err", err)
		return nil, err
	}
	if cfg.Service.Orchestration == utils.Executor {
		return executor.NewTraderService(ctx, &models.NewTraderServiceInput{
			ExchangeClient: exchangeClient,
			Trade:          models.TradeConfig(cfg.Trade),
			Executor:       models.ExecutorConfig(cfg.Executor),
			Logger:         logger,
		})
	} else if cfg.Service.Orchestration == utils.HTTP {
		return http.NewTraderService(ctx, &models.NewTraderServiceInput{
			ExchangeClient: exchangeClient,
			Trade:          models.TradeConfig(cfg.Trade),
			Executor:       models.ExecutorConfig(cfg.Executor),
			Logger:         logger,
		})
	} else {
		level.Error(logger).Log("msg", "invalid orchestration", "orchestration", cfg.Service.Orchestration)
		return nil, fmt.Errorf("invalid orchestration")
	}
}
