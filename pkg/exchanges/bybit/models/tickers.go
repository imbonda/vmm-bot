package models

import "fmt"

type RawTickersResult struct {
	Category     string      `json:"category"`
	List         []RawTicker `json:"list"`
	latestTicker *RawTicker
}

type RawTicker struct {
	Symbol                 string `json:"symbol"`
	LastPrice              string `json:"lastPrice"`
	IndexPrice             string `json:"indexPrice"`
	MarkPrice              string `json:"markPrice"`
	PrevPrice24h           string `json:"prevPrice24h"`
	Price24hPcnt           string `json:"price24hPcnt"`
	HighPrice24h           string `json:"highPrice24h"`
	LowPrice24h            string `json:"lowPrice24h"`
	PrevPrice1h            string `json:"prevPrice1h"`
	OpenInterest           string `json:"openInterest"`
	OpenInterestValue      string `json:"openInterestValue"`
	Turnover24h            string `json:"turnover24h"`
	Volume24h              string `json:"volume24h"`
	FundingRate            string `json:"fundingRate"`
	NextFundingTime        string `json:"nextFundingTime"`
	PredictedDeliveryPrice string `json:"predictedDeliveryPrice"`
	BasisRate              string `json:"basisRate"`
	DeliveryFeeRate        string `json:"deliveryFeeRate"`
	DeliveryTime           string `json:"deliveryTime"`
	BestAskPrice           string `json:"ask1Price"`
	BestAskQty             string `json:"ask1Size"`
	BestBidPrice           string `json:"bid1Price"`
	BestBidQty             string `json:"bid1Size"`
	Basis                  string `json:"basis"`
}

func (r *RawTickersResult) LatestTicker() (*RawTicker, error) {
	if r.latestTicker != nil {
		return r.latestTicker, nil
	}
	if len(r.List) < 1 {
		return nil, fmt.Errorf("no tickers found")
	}
	r.latestTicker = &r.List[0]
	return r.latestTicker, nil
}
