package fn

type Filter[T any] func(T) bool

type FilterWithError[T any] func(T) (bool, error)
