package fn

type Map[T any, N any] func(T) N

type MapWithError[T any, N any] func(T) (N, error)
