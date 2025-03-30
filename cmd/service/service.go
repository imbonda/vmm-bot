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
	"github.com/imbonda/bybit-vmm-bot/pkg/exchanges/biconomy"
	"github.com/imbonda/bybit-vmm-bot/pkg/utils"
)

func GetTraderService(ctx context.Context, cfg *config.Configuration) (interfaces.TraderService, error) {
	logger := cfg.GetLogger()
	bybitClient, err := biconomy.NewClient(ctx, &biconomy.NewClientInput{
		APIKey:    cfg.ExchangeAPIKey,
		APISecret: cfg.ExchangeAPISecret,
		Logger:    nil,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to create bybit client", "err", err)
		return nil, err
	}
	if cfg.ServiceOrchestration == utils.Executor {
		return executor.NewTraderService(ctx, &models.NewTraderServiceInput{
			ExchangeClient:                 bybitClient,
			Symbol:                         cfg.Symbol,
			IntervalExecutionDuration:      cfg.IntervalExecutionDuration,
			NumOfTradeIterationsInInterval: cfg.NumOfTradeIterationsInInterval,
			Logger:                         logger,
		})
	} else if cfg.ServiceOrchestration == utils.HTTP {
		return http.NewTraderService(ctx, &models.NewTraderServiceInput{
			ExchangeClient: bybitClient,
			Symbol:         cfg.Symbol,
			ListenAddress:  cfg.ListenAddress,
			Logger:         logger,
		})
	} else {
		level.Error(logger).Log("msg", "invalid orchestration", "orchestration", cfg.ServiceOrchestration)
		return nil, fmt.Errorf("invalid orchestration")
	}
}
