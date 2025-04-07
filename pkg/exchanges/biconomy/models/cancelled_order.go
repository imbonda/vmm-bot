package models

type CancelledOrderParam struct {
	Symbol  string `json:"market"`
	OrderId int    `json:"order_id"`
}

type CancelledOrder PendingOrder

type CancelledBatch []CancelledOrderInBatch

type CancelledOrderInBatch struct {
	Symbol     string `json:"market"`
	OrderId    int    `json:"order_id"`
	Successful string `json:"result"`
}
