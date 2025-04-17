package models

type RawPendingOrder struct {
	Symbol        string `json:"symbol"`
	OrderID       int    `json:"orderId"`
	ClientOrderID string `json:"clientOrderID"`
	Timestamp     int    `json:"transactTime"`
	Price         string `json:"price"`
	OrigQty       string `json:"origQty"`
	ExecQty       string `json:"executedQty"`
	CummQuoteQty  string `json:"cummulativeQuoteQty"`
	Status        string `json:"status"`
	Type          string `json:"type"`
	Side          string `json:"side"`
}
