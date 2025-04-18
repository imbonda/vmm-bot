package interfaces

import (
	"context"

	"github.com/imbonda/vmm-bot/pkg/models"
)

type ExchangeClient interface {
	GetOrderBook(ctx context.Context, symbol string) (*models.OrderBook, error)
	GetLastTicker(ctx context.Context, symbol string) (*models.Ticker, error)
	PlaceOrder(ctx context.Context, order *models.Order) error
	CancelAllOrders(ctx context.Context, symbol string) error
}
