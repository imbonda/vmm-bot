// pkg/trader/trader.go
package trader

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/imbonda/bybit-vmm-bot/pkg/api"
	"github.com/imbonda/bybit-vmm-bot/pkg/models"
)

const (
	apiKey    = "YOUR_API_KEY"
	apiSecret = "YOUR_API_SECRET"
	symbol    = "BTCUSDT" // Change to your desired trading pair
)

type Trader struct {
	client *api.BybitClient
}

func NewTrader() *Trader {
	return &Trader{
		client: api.NewBybitClient(apiKey, apiSecret),
	}
}

func (t *Trader) Start() error {
	for {
		orderBook, err := t.client.GetOrderBook(symbol)
		if err != nil {
			return fmt.Errorf("error fetching order book: %w", err)
		}

		fmt.Println(orderBook)      // TODO: remove
		fmt.Println(orderBook.List) // TODO: remove

		orderPrice := 0.0 // TODO: rand.Float64()*(orderBook.Asks[0].Price-orderBook.Bids[0].Price) + orderBook.Bids[0].Price
		side := "Buy"
		if rand.Intn(2) == 0 {
			side = "Sell"
		}

		orderQty := rand.Float64() * 0.01 // Random quantity between 0 and 0.01
		order := models.Order{
			Symbol: symbol,
			Side:   side,
			Price:  orderPrice,
			Qty:    orderQty,
		}

		if _, err := t.client.PlaceOrder(
			order.Symbol,
			order.Side,
			"Limit",
			fmt.Sprintf("%.4f", order.Qty),
			fmt.Sprintf("%.2f", order.Price),
		); err != nil {
			fmt.Println("Error placing order:", err)
		} else {
			fmt.Printf("Placed %s order for %.4f at price %.2f\n", order.Side, order.Qty, order.Price)
		}

		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second) // Sleep between 1 to 5 seconds
	}
}
