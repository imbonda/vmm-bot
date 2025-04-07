package models

import (
	"fmt"

	bybit "github.com/bybit-exchange/bybit.go.api"
)

const (
	successCode int = 0
)

type Response bybit.ServerResponse

func (r *Response) Validate() error {
	if !r.IsSuccessful() {
		err := fmt.Errorf("bybit request failed: %v", r.RetMsg)
		return err
	}
	return nil
}

func (r *Response) IsSuccessful() bool {
	return r.RetCode == successCode
}
