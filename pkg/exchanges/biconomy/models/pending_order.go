package models

type PendingOrdersResult struct {
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
	Records []RawPendingOrder `json:"records"`
}

type RawPendingOrder struct {
	Amount     string  `json:"amount"`
	CreatedAt  float64 `json:"ctime"`
	DealFee    string  `json:"deal_fee"`
	DealMoney  string  `json:"deal_money"`
	DealStock  string  `json:"deal_stock"`
	OrderId    int     `json:"id"`
	Left       string  `json:"left"`
	MakerFee   string  `json:"maker_fee"`
	Symbol     string  `json:"market"`
	ModifiedAt float64 `json:"mtime"`
	Price      string  `json:"price"`
	Side       int     `json:"side"`
	Source     string  `json:"source"`
	TakerFee   string  `json:"taker_fee"`
	Type       int     `json:"type"`
	User       int     `json:"user"`
}
