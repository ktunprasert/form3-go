package domain

type Request[T any] struct {
	Data *T `json:"data,omitempty"`
}
