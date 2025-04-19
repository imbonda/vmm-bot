package bingx

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-resty/resty/v2"

	"github.com/imbonda/vmm-bot/pkg/exchanges/bingx/hooks"
	bingxModels "github.com/imbonda/vmm-bot/pkg/exchanges/bingx/models"
	"github.com/imbonda/vmm-bot/pkg/models"
	"github.com/imbonda/vmm-bot/pkg/utils"
)

// API Configuration
const (
	BaseAPIURL = "https://open-api.bingx.com"
	APIV1      = "openApi/spot/v1"
)

type tradingSide string

const (
	SELL = "SELL"
	BUY  = "BUY"
)

func resolveSide(action models.OrderAction) tradingSide {
	if action == models.Buy {
		return BUY
	}
	return SELL
}

type Client struct {
	v1     *utils.Endpoint
	creds  *utils.Credentials
	client *resty.Client
	logger log.Logger
}

type NewClientInput struct {
	APIKey     string
	APISecret  string
	APITimeout time.Duration
	Logger     log.Logger
}

func NewClient(ctx context.Context, input *NewClientInput) (*Client, error) {
	v1 := utils.NewEndpoint(APIV1)
	creds := &utils.Credentials{
		APIKey:    input.APIKey,
		APISecret: input.APISecret,
	}
	client := resty.New().
		SetBaseURL(BaseAPIURL).
		SetHeader("Content-Type", "application/json").
		SetTimeout(input.APITimeout)
	// Add credentials to every request.
	client.OnBeforeRequest(hooks.GetSigAuthBeforeRequestHook(client, creds))
	return &Client{
		v1:     v1,
		creds:  creds,
		client: client,
		logger: input.Logger,
	}, nil
}

func (api *Client) GetOrderBook(ctx context.Context, symbol string) (*models.OrderBook, error) {
	bookTicker, err := api.getOrderBookTicker(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return &models.OrderBook{
		Symbol: symbol,
		Asks: [][]string{
			{bookTicker.AskPrice, bookTicker.AskAmount},
		},
		Bids: [][]string{
			{bookTicker.BidPrice, bookTicker.BidAmount},
		},
	}, nil
}

func (api *Client) GetLastTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	var bookTicker *bingxModels.BookTicker
	var priceTicker *bingxModels.PriceTicker
	var err1, err2 error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		bookTicker, err1 = api.getOrderBookTicker(ctx, symbol)
	}()

	go func() {
		defer wg.Done()
		priceTicker, err2 = api.getPriceTicker(ctx, symbol)
	}()

	wg.Wait()

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	return &models.Ticker{
		Symbol:    symbol,
		LastPrice: priceTicker.LastPrice,
		BestAsk:   bookTicker.AskPrice,
		BestBid:   bookTicker.BidPrice,
	}, nil
}

func (api *Client) getOrderBookTicker(ctx context.Context, symbol string) (*bingxModels.BookTicker, error) {
	var res bingxModels.Response[bingxModels.RawBookTickers]
	resp, err := api.client.R().
		SetResult(&res).
		SetQueryParam("symbol", symbol).
		Get(api.v1.Join("ticker/bookTicker"))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("bingx orderbook request failed with status: %s", resp.Status())
	}
	if !res.IsSuccessful() {
		return nil, fmt.Errorf("bingx orderbook request failed: %s", res.Message)
	}
	bookTicker, err := res.Result.LastTicker()
	if err != nil {
		return nil, err
	}
	return &bingxModels.BookTicker{
		AskPrice:  bookTicker.AskPrice,
		AskAmount: bookTicker.AskAmount,
		BidPrice:  bookTicker.BidPrice,
		BidAmount: bookTicker.BidAmount,
	}, nil
}

func (api *Client) getPriceTicker(ctx context.Context, symbol string) (*bingxModels.PriceTicker, error) {
	var res bingxModels.Response[bingxModels.RawPriceTickers]
	resp, err := api.client.R().
		SetResult(&res).
		SetQueryParam("symbol", symbol).
		Get(api.v1.Join("ticker/price"))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("bingx orderbook request failed with status: %s", resp.Status())
	}
	if !res.IsSuccessful() {
		return nil, fmt.Errorf("bingx orderbook request failed: %s", res.Message)
	}
	priceTicker, err := res.Result.LastTicker()
	if err != nil {
		return nil, err
	}
	lastTrade, err := priceTicker.LastTrade()
	if err != nil {
		return nil, err
	}
	return &bingxModels.PriceTicker{
		LastPrice: lastTrade.Price,
	}, nil
}

func (api *Client) PlaceOrder(ctx context.Context, order *models.Order) error {
	var res bingxModels.Response[bingxModels.RawPendingOrder]

	formData := map[string]string{
		"type":     "LIMIT",
		"symbol":   order.Symbol,
		"side":     string(resolveSide(order.Action)),
		"quantity": order.Qty,
		"price":    order.Price,
	}

	resp, err := api.client.R().
		SetFormData(formData).
		SetResult(&res).
		Post(api.v1.Join("trade/order"))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("bingx placeOrder request failed with status: %s", resp.Status())
	}
	if !res.IsSuccessful() {
		return fmt.Errorf("bingx placeOrder request failed: %s", res.Message)
	}

	return nil
}

func (api *Client) CancelAllOrders(ctx context.Context, symbol string) error {
	var res bingxModels.Response[bingxModels.RawCancelledBatch]

	formData := map[string]string{
		"symbol": symbol,
	}

	resp, err := api.client.R().
		SetFormData(formData).
		SetResult(&res).
		Post(api.v1.Join("trade/cancelOpenOrders"))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("bingx cancelAllOrders request failed with status: %s", resp.Status())
	}
	if !res.IsSuccessful() {
		return fmt.Errorf("bingx cancelAllOrders request failed: %s", res.Message)
	}

	return nil
}
