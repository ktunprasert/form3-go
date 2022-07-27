package domain

type Response[T any] struct {
	Data *T `json:"data,omitempty"`
}
