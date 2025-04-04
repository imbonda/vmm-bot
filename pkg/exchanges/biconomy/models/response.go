package models

const (
	successCode int = 0
)

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  T      `json:"result"`
}

func (r *Response[T]) IsSuccessful() bool {
	return r.Code == successCode
}
