package models

import (
	"fmt"

	"github.com/samber/lo"
)

type RawTickersResult struct {
	Tickers    []RawTicker `json:"ticker"`
	lastTicker *RawTicker
}

type RawTicker struct {
	Symbol    string `json:"symbol"`
	Change    string `json:"change"`
	Deal      string `json:"deal"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Volume24h string `json:"vol"`
	LastPrice string `json:"last"`
	Ask       string `json:"sell"`
	Bid       string `json:"buy"`
}

func (r *RawTickersResult) LastTicker(symbol string) (*RawTicker, error) {
	if r.lastTicker != nil {
		return r.lastTicker, nil
	}
	err := fmt.Errorf("no tickers found")
	if len(r.Tickers) < 1 {
		return nil, err
	}
	ticker, found := lo.Find(
		r.Tickers,
		func(ticker RawTicker) bool {
			return ticker.Symbol == symbol
		},
	)
	if !found {
		return nil, err
	}
	r.lastTicker = &ticker
	return r.lastTicker, nil
}
