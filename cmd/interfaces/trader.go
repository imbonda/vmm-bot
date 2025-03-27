package interfaces

import (
	"context"

	"github.com/imbonda/bybit-vmm-bot/pkg/models"
)

type Trader interface {
	TradeOnce(ctx context.Context) (*models.TradeOnceOutput, error)
}
