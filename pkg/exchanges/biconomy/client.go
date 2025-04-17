package biconomy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"

	"github.com/imbonda/vmm-bot/pkg/exchanges/biconomy/hooks"
	biconomyModels "github.com/imbonda/vmm-bot/pkg/exchanges/biconomy/models"
	"github.com/imbonda/vmm-bot/pkg/models"
	"github.com/imbonda/vmm-bot/pkg/utils"
)

// API Configuration
const (
	BaseAPIURL = "https://market.biconomy.vip"
	APIV1      = "api/v1"
	APIV2      = "api/v2"
)

type tradingSide string

const (
	ASK = "1"
	BID = "2"
)

func resolveSide(action models.OrderAction) tradingSide {
	if action == models.Buy {
		return BID
	}
	return ASK
}

type Client struct {
	v1     *utils.Endpoint
	v2     *utils.Endpoint
	creds  *utils.Credentials
	client *resty.Client
	logger log.Logger
}

type NewClientInput struct {
	APIKey    string
	APISecret string
	Logger    log.Logger
}

func NewClient(ctx context.Context, input *NewClientInput) (*Client, error) {
	v1 := utils.NewEndpoint(APIV1)
	v2 := utils.NewEndpoint(APIV2)
	creds := &utils.Credentials{
		APIKey:    input.APIKey,
		APISecret: input.APISecret,
	}
	client := resty.New().
		SetBaseURL(BaseAPIURL).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("X-SITE-ID", "127")
	// Add credentials to every request.
	client.OnBeforeRequest(hooks.GetSigAuthBeforeRequestHook(client, creds))
	return &Client{
		v1:     v1,
		v2:     v2,
		creds:  creds,
		client: client,
		logger: input.Logger,
	}, nil
}

func (api *Client) GetOrderBook(ctx context.Context, symbol string) (*models.OrderBook, error) {
	var res biconomyModels.RawOrderBook
	resp, err := api.client.R().
		SetResult(&res).
		SetQueryParam("symbol", symbol).
		Get(api.v1.Join("depth"))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("biconomy orderbook request failed: %s", resp.Status())
	}
	return &models.OrderBook{
		Symbol: symbol,
		Asks:   res.Asks,
		Bids:   res.Bids,
	}, nil
}

func (api *Client) GetLastTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	var res biconomyModels.RawTickersResult
	resp, err := api.client.R().
		SetResult(&res).
		Get(api.v1.Join("tickers"))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("biconomy tickers request failed: %s", resp.Status())
	}
	ticker, err := res.LastTicker(symbol)
	if err != nil {
		return nil, err
	}
	return &models.Ticker{
		Symbol:    ticker.Symbol,
		LastPrice: ticker.LastPrice,
		BestAsk:   ticker.Ask,
		BestBid:   ticker.Bid,
	}, nil
}

func (api *Client) PlaceOrder(ctx context.Context, order *models.Order) error {
	var res biconomyModels.Response[biconomyModels.RawFulfilledOrder]

	formData := map[string]string{
		"market": order.Symbol,
		"amount": order.Qty,
		"price":  order.Price,
		"side":   string(resolveSide(order.Action)),
	}

	resp, err := api.client.R().
		SetFormData(formData).
		SetResult(&res).
		Post(api.v1.Join("private/trade/limit"))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("biconomy placeOrder request failed with status: %s", resp.Status())
	}
	if !res.IsSuccessful() {
		return fmt.Errorf("biconomy placeOrder request failed: %s", res.Message)
	}

	return nil
}

func (api *Client) CancelAllOrders(ctx context.Context, symbol string) error {
	records, err := api.queryUnfilledOrders(ctx, symbol)
	if err != nil {
		return err
	}
	if len(records) > 0 {
		return api.batchCancelOrders(ctx, records)
	}
	return nil
}

func (api *Client) queryUnfilledOrders(_ context.Context, symbol string) ([]biconomyModels.PendingOrder, error) {
	var res biconomyModels.Response[biconomyModels.PendingOrdersResult]

	formData := map[string]string{
		"market": symbol,
		"limit":  "10",
	}

	resp, err := api.client.R().
		SetFormData(formData).
		SetResult(&res).
		Post(api.v1.Join("private/order/pending"))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("biconomy queryUnfilledOrders request failed with status: %s", resp.Status())
	}
	if !res.IsSuccessful() {
		return nil, fmt.Errorf("biconomy queryUnfilledOrders request failed: %s", res.Message)
	}

	return res.Result.Records, nil
}

func (api *Client) cancelOrder(_ context.Context, order *biconomyModels.PendingOrder) error {
	var res biconomyModels.Response[biconomyModels.CancelledOrder]

	formData := map[string]string{
		"market":   order.Symbol,
		"order_id": utils.FormatIntToString(order.OrderId),
	}

	resp, err := api.client.R().
		SetFormData(formData).
		SetResult(&res).
		Post(api.v1.Join("private/trade/cancel"))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("biconomy cancelOrder request failed with status: %s", resp.Status())
	}
	if !res.IsSuccessful() {
		return fmt.Errorf("biconomy cancelOrder request failed: %s", res.Message)
	}

	return nil
}

func (api *Client) batchCancelOrders(_ context.Context, orders []biconomyModels.PendingOrder) error {
	var res biconomyModels.Response[biconomyModels.CancelledBatch]

	ordersParams := lo.Map(orders, func(order biconomyModels.PendingOrder, _ int) biconomyModels.CancelledOrderParam {
		return biconomyModels.CancelledOrderParam{
			Symbol:  order.Symbol,
			OrderId: order.OrderId,
		}
	})
	ordersJson, err := json.Marshal(ordersParams)
	if err != nil {
		return err
	}
	formData := map[string]string{
		"orders_json": string(ordersJson),
	}

	resp, err := api.client.R().
		SetFormData(formData).
		SetResult(&res).
		Post(api.v1.Join("private/trade/cancel_batch"))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("biconomy batchCancelOrders request failed with status: %s", resp.Status())
	}
	if !res.IsSuccessful() {
		return fmt.Errorf("biconomy batchCancelOrders request failed: %s", res.Message)
	}

	return nil
}
