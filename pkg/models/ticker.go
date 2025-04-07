package models

import (
	"fmt"

	"github.com/imbonda/vmm-bot/pkg/utils"
)

type Ticker struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"lastPrice"`
	BestAsk   string `json:"ask"`
	BestBid   string `json:"bid"`
}

func (t *Ticker) Spread() (*Spread, error) {
	return NewSpread(t.BestAsk, t.BestBid)
}

func (t *Ticker) Price() (float64, error) {
	price, err := utils.ParseFloat(t.LastPrice)
	if err != nil {
		return 0, fmt.Errorf("failed parse last price: %s", t.LastPrice)
	}
	return price, nil
}
