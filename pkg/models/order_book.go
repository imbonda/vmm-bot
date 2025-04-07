package models

import (
	"fmt"
)

type OrderBook struct {
	Symbol  string     `json:"s"`
	Asks    [][]string `json:"a"`
	Bids    [][]string `json:"b"`
	bestAsk *string
	bestBid *string
}

func (b *OrderBook) Ask() (string, error) {
	if b.bestAsk != nil {
		return *b.bestAsk, nil
	}
	if len(b.Asks) == 0 {
		return "", nil
	}
	if len(b.Asks[0]) < 2 {
		return "", fmt.Errorf("invalid order book")
	}
	askPrice := b.Asks[0][0]
	b.bestAsk = &askPrice
	return askPrice, nil
}

func (b *OrderBook) Bid() (string, error) {
	if b.bestBid != nil {
		return *b.bestBid, nil
	}
	if len(b.Bids) == 0 {
		return "", nil
	}
	if len(b.Bids[0]) < 2 {
		return "", fmt.Errorf("invalid order book")
	}
	bidPrice := b.Bids[0][0]
	b.bestBid = &bidPrice
	return bidPrice, nil
}

func (b *OrderBook) Spread() (*Spread, error) {
	ask, err := b.Ask()
	if err != nil {
		return nil, err
	}
	bid, err := b.Bid()
	if err != nil {
		return nil, err
	}
	return NewSpread(ask, bid)
}
