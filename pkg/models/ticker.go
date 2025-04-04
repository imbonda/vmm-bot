package models

type Ticker struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"lastPrice"`
	BestAsk   string `json:"ask1Price"`
	BestBid   string `json:"bid1Price"`
}

func (t *Ticker) Spread() (*Spread, error) {
	return NewSpread(t.BestAsk, t.BestBid)
}
