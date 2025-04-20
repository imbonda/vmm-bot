package bybit

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

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
	APIKey     string
	APISecret  string
	APITimeout time.Duration
	Logger     log.Logger
}

func NewClient(ctx context.Context, input *NewClientInput) (*Client, error) {
	return &Client{
		client: bybit.NewBybitHttpClient(
			input.APIKey,
			input.APISecret,
			bybit.WithBaseURL(bybit.MAINNET),
			func(c *bybit.Client) {
				c.HTTPClient.Timeout = input.APITimeout
			},
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
	err = wrapResponseErrors(
		&requestInfo{
			method:   http.MethodGet,
			endpoint: "GetOrderBookInfo",
		},
		err,
	)
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

func (api *Client) GetLastTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	res, err := api.client.
		NewUtaBybitServiceWithParams(
			map[string]any{
				"category": "spot",
				"symbol":   symbol,
			},
		).
		GetMarketTickers(ctx)
	err = wrapResponseErrors(
		&requestInfo{
			method:   http.MethodGet,
			endpoint: "GetMarketTickers",
		},
		err,
	)
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
	ticker, err := rawResult.LastTicker()
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
	err = wrapResponseErrors(
		&requestInfo{
			method:   http.MethodPost,
			endpoint: "PlaceOrder",
		},
		err,
	)
	if err != nil {
		return err
	}
	wrappedRes := bybitModels.Response(*res)
	if err = wrappedRes.Validate(); err != nil {
		return err
	}
	return nil
}

func (api *Client) CancelAllOrders(ctx context.Context, symbol string) error {
	res, err := api.client.
		NewUtaBybitServiceWithParams(
			map[string]any{
				"category": "spot",
				"symbol":   symbol,
			},
		).
		CancelAllOrders(context.Background())
	err = wrapResponseErrors(
		&requestInfo{
			method:   http.MethodPost,
			endpoint: "CancelAllOrders",
		},
		err,
	)
	if err != nil {
		return err
	}
	wrappedRes := bybitModels.Response(*res)
	if err = wrappedRes.Validate(); err != nil {
		return err
	}
	return nil
}
