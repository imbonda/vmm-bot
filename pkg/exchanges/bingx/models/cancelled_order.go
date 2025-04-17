package models

type RawCancelledOrder RawPendingOrder

type RawCancelledBatch struct {
	Orders []RawPendingOrder `json:"orders"`
}
