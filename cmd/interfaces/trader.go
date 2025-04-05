package interfaces

import (
	"context"

	"github.com/imbonda/vmm-bot/pkg/models"
)

type Trader interface {
	TradeOnce(ctx context.Context) (*models.TradeOnceOutput, error)
}
