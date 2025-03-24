package interfaces

import (
	"context"

	"github.com/imbonda/bybit-vmm-bot/pkg/models"
)

type ExchangeClient interface {
	GetOrderBook(ctx context.Context, symbol string) (*models.OrderBook, error)
	PlaceOrder(ctx context.Context, order *models.Order) error
}
