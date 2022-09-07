package domain

// TODO: Add pagination
type List[T any] struct {
	Data []T `json:"data"`
}

func NewList[T any](data []T) *List[T] {
	return &List[T]{
		Data: data,
	}
}
