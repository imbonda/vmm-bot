package bybit

import (
	"context"
	"encoding/json"

	bybit "github.com/bybit-exchange/bybit.go.api"
	"github.com/go-kit/log"

	bybitModels "github.com/imbonda/vmm-bot/pkg/exchanges/bybit/models"
	"github.com/imbonda/vmm-bot/pkg/models"
)

type Client struct {
	client *bybit.Client
	logger log.Logger
}

type NewClientInput struct {
	APIKey    string
	APISecret string
	Logger    log.Logger
}

func NewClient(ctx context.Context, input *NewClientInput) (*Client, error) {
	return &Client{
		client: bybit.NewBybitHttpClient(
			input.APIKey,
			input.APISecret,
			bybit.WithBaseURL(bybit.MAINNET),
		),
	}, nil
}

func (api *Client) GetOrderBook(ctx context.Context, symbol string) (*models.OrderBook, error) {
	res, err := api.client.
		NewUtaBybitServiceWithParams(
			map[string]any{
				"category": "spot",
				"symbol":   symbol,
			},
		).
		GetOrderBookInfo(ctx)
	if err != nil {
		return nil, err
	}
	wrappedRes := bybitModels.Response(*res)
	if err = wrappedRes.Validate(); err != nil {
		return nil, err
	}
	data, err := json.Marshal(res.Result)
	if err != nil {
		return nil, err
	}
	result := &models.OrderBook{}
	err = json.Unmarshal(data, result)
	return result, err
}

func (api *Client) GetLatestTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	res, err := api.client.
		NewUtaBybitServiceWithParams(
			map[string]any{
				"category": "spot",
				"symbol":   symbol,
			},
		).
		GetMarketTickers(ctx)
	if err != nil {
		return nil, err
	}
	wrappedRes := bybitModels.Response(*res)
	if err = wrappedRes.Validate(); err != nil {
		return nil, err
	}
	data, err := json.Marshal(res.Result)
	if err != nil {
		return nil, err
	}
	rawResult := &bybitModels.RawTickersResult{}
	err = json.Unmarshal(data, rawResult)
	if err != nil {
		return nil, err
	}
	ticker, err := rawResult.LatestTicker()
	if err != nil {
		return nil, err
	}
	result := &models.Ticker{
		Symbol:    ticker.Symbol,
		LastPrice: ticker.LastPrice,
		BestAsk:   ticker.BestAskPrice,
		BestBid:   ticker.BestBidPrice,
	}
	return result, nil
}

func (api *Client) PlaceOrder(ctx context.Context, order *models.Order) error {
	res, err := api.client.
		NewUtaBybitServiceWithParams(
			map[string]any{
				"category":    "spot",
				"symbol":      order.Symbol,
				"side":        order.Action,
				"positionIdx": 0,
				"orderType":   "Limit",
				"qty":         order.Qty,
				"price":       order.Price,
				"timeInForce": "GTC",
			},
		).
		PlaceOrder(context.Background())
	if err != nil {
		return err
	}
	wrappedRes := bybitModels.Response(*res)
	if err = wrappedRes.Validate(); err != nil {
		return err
	}
	return nil
}
