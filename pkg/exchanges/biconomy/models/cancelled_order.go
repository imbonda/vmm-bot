package models

type CancelledOrderParam struct {
	Symbol  string `json:"market"`
	OrderId int    `json:"order_id"`
}

type RawCancelledOrder RawPendingOrder

type RawCancelledBatch []RawCancelledOrderInBatch

type RawCancelledOrderInBatch struct {
	Symbol     string `json:"market"`
	OrderId    int    `json:"order_id"`
	Successful string `json:"result"`
}
